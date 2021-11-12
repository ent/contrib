// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AttachmentsColumns holds the columns for the "attachments" table.
	AttachmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "user_attachment", Type: field.TypeInt, Unique: true, Nullable: true},
	}
	// AttachmentsTable holds the schema information for the "attachments" table.
	AttachmentsTable = &schema.Table{
		Name:       "attachments",
		Columns:    AttachmentsColumns,
		PrimaryKey: []*schema.Column{AttachmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "attachments_users_attachment",
				Columns:    []*schema.Column{AttachmentsColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// GroupsColumns holds the columns for the "groups" table.
	GroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
	}
	// GroupsTable holds the schema information for the "groups" table.
	GroupsTable = &schema.Table{
		Name:       "groups",
		Columns:    GroupsColumns,
		PrimaryKey: []*schema.Column{GroupsColumns[0]},
	}
	// NilExamplesColumns holds the columns for the "nil_examples" table.
	NilExamplesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "str_nil", Type: field.TypeString, Nullable: true},
		{Name: "time_nil", Type: field.TypeTime, Nullable: true},
	}
	// NilExamplesTable holds the schema information for the "nil_examples" table.
	NilExamplesTable = &schema.Table{
		Name:       "nil_examples",
		Columns:    NilExamplesColumns,
		PrimaryKey: []*schema.Column{NilExamplesColumns[0]},
	}
	// TodosColumns holds the columns for the "todos" table.
	TodosColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "task", Type: field.TypeString},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"pending", "in_progress", "done"}, Default: "pending"},
		{Name: "todo_user", Type: field.TypeInt, Nullable: true},
	}
	// TodosTable holds the schema information for the "todos" table.
	TodosTable = &schema.Table{
		Name:       "todos",
		Columns:    TodosColumns,
		PrimaryKey: []*schema.Column{TodosColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "todos_users_user",
				Columns:    []*schema.Column{TodosColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "user_name", Type: field.TypeString, Unique: true},
		{Name: "joined", Type: field.TypeTime},
		{Name: "points", Type: field.TypeUint},
		{Name: "exp", Type: field.TypeUint64},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"pending", "active"}},
		{Name: "external_id", Type: field.TypeInt, Unique: true},
		{Name: "crm_id", Type: field.TypeUUID},
		{Name: "banned", Type: field.TypeBool, Default: false},
		{Name: "custom_pb", Type: field.TypeUint8},
		{Name: "opt_num", Type: field.TypeInt, Nullable: true},
		{Name: "opt_str", Type: field.TypeString, Nullable: true},
		{Name: "opt_bool", Type: field.TypeBool, Nullable: true},
		{Name: "opt_strings", Type: field.TypeJSON, Nullable: true},
		{Name: "big_int", Type: field.TypeInt, Nullable: true},
		{Name: "b_user_1", Type: field.TypeInt, Unique: true, Nullable: true},
		{Name: "height_in_cm", Type: field.TypeFloat32, Default: 0},
		{Name: "account_balance", Type: field.TypeFloat64, Default: 0},
		{Name: "strings", Type: field.TypeJSON},
		{Name: "user_group", Type: field.TypeInt, Nullable: true},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "users_groups_group",
				Columns:    []*schema.Column{UsersColumns[19]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// AttachmentRecipientsColumns holds the columns for the "attachment_recipients" table.
	AttachmentRecipientsColumns = []*schema.Column{
		{Name: "attachment_id", Type: field.TypeUUID},
		{Name: "user_id", Type: field.TypeInt},
	}
	// AttachmentRecipientsTable holds the schema information for the "attachment_recipients" table.
	AttachmentRecipientsTable = &schema.Table{
		Name:       "attachment_recipients",
		Columns:    AttachmentRecipientsColumns,
		PrimaryKey: []*schema.Column{AttachmentRecipientsColumns[0], AttachmentRecipientsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "attachment_recipients_attachment_id",
				Columns:    []*schema.Column{AttachmentRecipientsColumns[0]},
				RefColumns: []*schema.Column{AttachmentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "attachment_recipients_user_id",
				Columns:    []*schema.Column{AttachmentRecipientsColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AttachmentsTable,
		GroupsTable,
		NilExamplesTable,
		TodosTable,
		UsersTable,
		AttachmentRecipientsTable,
	}
)

func init() {
	AttachmentsTable.ForeignKeys[0].RefTable = UsersTable
	TodosTable.ForeignKeys[0].RefTable = UsersTable
	UsersTable.ForeignKeys[0].RefTable = GroupsTable
	AttachmentRecipientsTable.ForeignKeys[0].RefTable = AttachmentsTable
	AttachmentRecipientsTable.ForeignKeys[1].RefTable = UsersTable
}
