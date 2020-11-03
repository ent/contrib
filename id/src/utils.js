import {camelCase, startCase, toUpper} from 'lodash';

class NullValueError extends Error {
  constructor(message?: string) {
    super('[NullValueError]' + (message ? ' ' + message : ''));
  }
}

export function nullthrows<TVal>(data: ?TVal, message?: string): TVal {
  if (data == null) {
    throw new NullValueError(message);
  }
  return data;
}

// http://github.com/golang/lint/blob/master/lint.go
const commonGoInitialisms = [
  'ACL',
  'API',
  'ASCII',
  'CPU',
  'CSS',
  'DNS',
  'EOF',
  'GUID',
  'HTML',
  'HTTP',
  'HTTPS',
  'ID',
  'IP',
  'JSON',
  'LHS',
  'QPS',
  'RAM',
  'RHS',
  'RPC',
  'SLA',
  'SMTP',
  'SQL',
  'SSH',
  'TCP',
  'TLS',
  'TTL',
  'UDP',
  'UI',
  'UID',
  'UUID',
  'URI',
  'URL',
  'UTF8',
  'VM',
  'XML',
  'XMPP',
  'XSRF',
  'XSS',
];

export const pascalCaseGoStyle = (word: string) => {
  return startCase(camelCase(word))
    .split(' ')
    .map(w => (commonGoInitialisms.includes(toUpper(w)) ? toUpper(w) : w))
    .join('');
};