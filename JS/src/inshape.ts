export class ShapeServa {
  Address: string;
  ws: WebSocket;
  private AuthTocken: string;
  constructor(Address: string) {
    this.Address = Address;
  }
  GetToken(): Promise<string> {
    return new Promise((resolve, reject) => {
      this.htMSG("GetToken", null)
        .then((res) => {
          console.log(res);
          console.log(res, "res");
          if (res !== "") {
            this.AuthTocken = res;
            resolve(this.AuthTocken);
          }
        })
        .catch(reject);
    });
  }
  Post(msg: any): Promise<any> {
    return new Promise((resolve, reject) => {
      fetch(this.Address, {
        method: "post",
        headers: { _Auth_Tocken_: this.AuthTocken },
        body: JSON.stringify(msg),
      })
        .then((res) => {
          console.log("re", res, msg);
          res
            .json()
            .then((res) => {
              if (res.Err === null) {
                // auth.Tocken = res.Tocken
                resolve(res.Data);
              } else {
                reject(res);
              }
            })
            .catch((err) => {
              reject(err);
            });
        })
        .catch((err) => reject(err));
    });
  }
  CreateShape(data: { Name: string; Methods: any[] }, Conf: any) {
    let obj: any = {};
    data.Methods.forEach((item, ind) => {
      obj[item] = (...args: any[]) => {
        return this.htMSG("Shape", [data.Name, item, args]);
      };
    });
    return obj;
  }

  private wsMSG(...args) {}
  private htMSG(Intent, Data): Promise<any> {
    return this.Post({ Intent, Data });
  }
  async RequestShape(Name: string, Conf: any) {
    let res = await this.htMSG("RequestShape", Name);
    return this.CreateShape(res, Conf);
  }
}

export class Shape {
  _sss: {
    Alive: boolean;
    Serva: ShapeServa;
    Name: string;
    InstanceID: string;
    // id: number;
  };
  _Init: any;
  _Events: { [key: string]: (...args: any[]) => void };
  Listeners: { [key: string]: any };
  constructor(Name: string, Serva: ShapeServa) {
    this._sss = {
      Alive: false,
      Serva: Serva,
      Name: Name,
      InstanceID: "",
    };
    this._Init = async (Conf?: any) => {
      let res = await this._sss.Serva.RequestShape(Name, Conf);
      Object.assign(this, res);
      this._sss.Alive = true;
      return this;
    };
    Object.defineProperty(this, "_Init", {
      configurable: false,
      writable: false,
    });
    Object.defineProperty(this, "_Events", {
      configurable: false,
      writable: false,
    });
  }
}
