import { Operation } from './Operation';

export class ListOperation implements Operation {
  parseOperation(): string {
    return 'op=list;\n';
  }
}
