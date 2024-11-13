# RBAC Auth

Ketika init Auth RBAC, semua host *required*.

## Menggunakan Auth RBAC
1. Init client RBAC
```Go
NewRBAC(
    validator, 
    WithHTTPHost(httpHost),
    WithGRPCHost(grpcHost),
    WithApplicationCode(applicationCode),
    )
```

2. Masukkan ke dalam middleware Auth
```Go
middlewareAuth := commonAuth.NewAuth(
		...
		commonAuth.WithRBACClient(rbacClient),
```

3. Gunakan salahsatu function dari RBAC. Contoh, menggunakan `IsRoleAllowed()`
```Go
    // Router.go

    // jika role user ada di dalam map, maka user authorized. Jika tidak, return error 401
    r.engine.GET("/health",
		r.common.GetAuthMiddleware().AuthRoleRBAC(map[string]bool{"Dealer Consultant": true}, "111"),
        
		r.deliveryRegistry.GetHealthCheck().Do)
```
