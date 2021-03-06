import net from 'net';

import { Operation } from './Operation';
import { ReadOperation } from './ReadOperation';
import { StoreOperation } from './StoreOperation';
import { DeleteOperation } from './DeleteOperation';
import { ListOperation } from './ListOperation';
import { KeysOperation } from './KeysOperation';
import { Response } from './Response';

export default class Client {
  private client: net.Socket;
  private host: string;
  private port: number;

  constructor(host: string, port: number) {
    this.client = new net.Socket();
    this.host = host;
    this.port = port;
  }

  async connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.client.connect({ host: this.host, port: this.port }, () => {
        resolve();
      });

      // if connect fail, will emmit a error event
      this.client.on('error', (err) => {
        reject(err);
      });
    });
  }

  async storeOperation(key: string, value: any) {
    const storeOp = new StoreOperation(key, value);
    const response = await this.sendMessage(storeOp);
    return response
  }

  async readOperation(key: string): Promise<Response> {
    const readOp = new ReadOperation(key);
    const response = await this.sendMessage(readOp);
    return response;
  }

  async deleteOperation(key: string): Promise<Response> {
    const deleteOp = new DeleteOperation(key);
    const response = await this.sendMessage(deleteOp);
    return response;
  }

  async listOperation(): Promise<Response> {
    const listOp = new ListOperation();
    const response = await this.sendMessage(listOp);
    return response;
  }

  async keysOperation(): Promise<Response> {
    const keysOp = new KeysOperation();
    const response = await this.sendMessage(keysOp);
    return response;
  }

  private async sendMessage(op: Operation): Promise<Response> {
    const msg = op.parseOperation();
    return new Promise((resolve, reject) => {
      this.client.write(msg, () => {
        this.client.on('data', (data) => {
          const response = this.getResponse(data.toString()); 
          resolve(response);
        });
        this.client.on('error', (err) => {
          reject(err);
        });
      });
    });
  }

  private getResponse(data: string): Response {
    const response = Response.parseResponse(data.toString());
    return response;
  }

  disconnect() {
    this.client.destroy();
  }
}
