# zmidd
some gin middleware

### Usage
#### Cors
```
func main() {
    s := gin.New()
    // default: *
    s.Use(zmidd.Cors())
    
    // options
    s.Use(zmidd.Cors(
        zmidd.WithAllowOrigin("*"),
        zmidd.WithAllowMethods("*"),
        zmidd.WithAllowHeaders("*"),
        zmidd.WithAllowCredentials("*"),
        zmidd.WithExposeHeaders("true"),
    ))
}
```

#### RequestId
```
func main() {
    s := gin.New()
    s.Use(zmidd.RequestId())
    
    // Get RequestId
    // c = *gin.Context
    rid := zmidd.GetRequestId(c)
}
```