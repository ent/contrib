// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entprom

import (
	"context"
	"time"

	"entgo.io/ent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	mutationType = "mutation_type"
	mutationOp   = "mutation_op"
)

var entLabels = []string{mutationType, mutationOp}

func (p *hook) initOpsProcessedTotal() {
	p.opsProcessedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "ent_operation_total",
			Help:        "Number of ent mutation operations",
			ConstLabels: p.extraLabels,
		},
		entLabels,
	)
}

func (p *hook) initOpsProcessedError() {
	p.opsProcessedError = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "ent_operation_error",
			Help:        "Number of failed ent mutation operations",
			ConstLabels: p.extraLabels,
		},
		entLabels,
	)
}

func (p *hook) initOpsDuration() {
	p.opsDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "ent_operation_duration_seconds",
			Help:        "Time in seconds per operation",
			ConstLabels: p.extraLabels,
		},
		entLabels,
	)
}

func newHook(p *hook) ent.Hook {
	p.initOpsProcessedTotal()
	p.initOpsProcessedError()
	p.initOpsDuration()
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			start := time.Now()
			labels := prometheus.Labels{mutationType: m.Type(), mutationOp: m.Op().String()}
			p.opsProcessedTotal.With(labels).Inc()
			v, err := next.Mutate(ctx, m)
			if err != nil {
				p.opsProcessedError.With(labels).Inc()
			}
			duration := time.Since(start)
			p.opsDuration.With(labels).Observe(duration.Seconds())
			return v, err
		})
	}
}

type hook struct {
	extraLabels       map[string]string
	opsProcessedTotal *prometheus.CounterVec
	opsProcessedError *prometheus.CounterVec
	opsDuration       *prometheus.HistogramVec
}

type Option func(h *hook)

// Labels allows you to add constant labels to your metrics.
func Labels(labels map[string]string) Option {
	return func(h *hook) {
		h.extraLabels = labels
	}
}

// Hook sends ent metrics to your prometheus.
func Hook(options ...Option) ent.Hook {
	ph := &hook{}
	for _, option := range options {
		option(ph)
	}
	return newHook(ph)
}
