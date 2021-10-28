import { Operation } from './Operation';

export class KeysOperation implements Operation {
  parseOperation(): string {
    return 'op=keys;\n';
  }
}
