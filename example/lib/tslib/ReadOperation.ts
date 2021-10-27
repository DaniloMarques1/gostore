import { Operation } from './Operation';

export class ReadOperation implements Operation {
  private key: string;

  constructor(key: string) {
    this.key = key;
  }

  parseOperation(): string {
    return `op=read;key=${this.key};\n`;
  }
}
