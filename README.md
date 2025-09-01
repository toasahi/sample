sequenceDiagram
participant User as ユーザー
participant Client as クライアント<br/>アプリケーション
participant Keycloak as Keycloak<br/>認証サーバー
participant API as APIサーバー<br/>(リソースサーバー)
participant JWKS as JWKS<br/>エンドポイント

```
Note over User, JWKS: Authorization Code Flow with PKCE

%% 1. 初期アクセス
User->>Client: 1. アプリケーションにアクセス
Client->>Client: 2. 認証状態チェック<br/>(トークンなし)

%% 3. 認証開始
Client->>Client: 3. PKCE用のcode_verifier<br/>とcode_challenge生成
Client->>User: 4. Keycloakログインページへリダイレクト<br/>(client_id, redirect_uri, code_challenge含む)

%% 5. ユーザー認証
User->>Keycloak: 5. ログイン情報入力<br/>(ユーザー名/パスワード)
Keycloak->>Keycloak: 6. 認証情報検証

%% 7. 認可コード発行
Keycloak->>User: 7. 認可コードと共に<br/>クライアントにリダイレクト
User->>Client: 8. 認可コード受信<br/>(callback URL経由)

%% 9. トークン交換
Client->>Keycloak: 9. トークンリクエスト<br/>(code, code_verifier, client_id)
Keycloak->>Keycloak: 10. code_verifierとcode_challenge検証
Keycloak->>Client: 11. JWT Access Token<br/>+ Refresh Token返却

%% 12. API呼び出し
Client->>API: 12. APIリクエスト<br/>(Authorization: Bearer {JWT})

%% 13. JWT検証プロセス
API->>API: 13. JWTヘッダーから<br/>kid (Key ID) 抽出
API->>JWKS: 14. 公開鍵取得<br/>GET /auth/realms/{realm}/protocol/openid_connect/certs
JWKS->>API: 15. JWKSレスポンス<br/>(公開鍵セット)

API->>API: 16. JWT署名検証<br/>• kidで対応する公開鍵選択<br/>• RSA/ECDSA署名検証<br/>• exp, iat, iss検証

alt JWT検証成功
    API->>API: 17a. ペイロード抽出<br/>• ユーザー情報<br/>• 権限情報<br/>• 有効期限
    API->>Client: 18a. APIレスポンス<br/>(リクエストされたデータ)
    Client->>User: 19a. 結果表示
else JWT検証失敗
    API->>Client: 17b. 401 Unauthorized<br/>エラーレスポンス
    Client->>Client: 18b. トークンリフレッシュ<br/>または再認証処理
end

%% トークンリフレッシュフロー
Note over Client, Keycloak: Token Refresh Flow

Client->>Keycloak: 20. リフレッシュトークンで<br/>新しいアクセストークン要求
Keycloak->>Keycloak: 21. リフレッシュトークン検証
Keycloak->>Client: 22. 新しいJWT Access Token

%% ログアウト
Note over User, Keycloak: Logout Flow

User->>Client: 23. ログアウト要求
Client->>Keycloak: 24. ログアウトエンドポイント呼び出し<br/>(トークン無効化)
Keycloak->>Client: 25. ログアウト完了
Client->>Client: 26. ローカルトークン削除
Client->>User: 27. ログアウト完了通知
```