import { Operation } from './Operation';

export class DeleteOperation implements Operation {
  private key: string;

  constructor(key: string) {
    this.key = key;
  }

  parseOperation(): string {
    return `op=delete;key=${this.key};\n`;
  }
}
