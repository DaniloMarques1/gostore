import { Operation } from './Operation';

export class StoreOperation implements Operation {
  private key: string;
  private value: any;

  constructor(key: string, value: any) {
    if (typeof value === 'object') {
      value = JSON.stringify(value);
    }
    this.key = key;
    this.value = value;
  }

  parseOperation(): string {
    return `op=store;key=${this.key};value=${this.value};\n`;
  }
}
