export class Response {
  private code: number;
  private message: string;

  private constructor(code: number, message: string) {
    this.code = code;
    this.message = message;
  }
  
  static parseResponse(msg: string) {
    msg = msg.replace('\n', '');
    const splited = msg.split(';');
    const codeSplit = splited[0].split('=');
    const code = codeSplit[1];
    const messageSplit = splited[1].split('=');
    const message = messageSplit[1];

    return new Response(Number(code), message);
  }

  getCode(): number {
    return this.code;
  }

  getMessage(): string {
    return this.message;
  }

}
