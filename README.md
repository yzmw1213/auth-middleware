# auth-middleware
![architect](https://github.com/yzmw1213/auth-middleware/assets/36359899/2bf7b734-ae60-4d43-9f52-41068129091d)

## 概要
- Authorization HeaderにBearer tokenを設定し、HTTPリクエスト
- tokenを取得し、Firebaseに検証を投げる
- 返却されたClaimを検証。 
- user_idに該当するユーザーが存在する・かつ APIに対する権限条件が満たされる場合にロジックを実行する
