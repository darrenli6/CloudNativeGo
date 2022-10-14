module.exports = {
    devServer:{
        overlay: {
            warnings: true,
            errors: true
        },
        host: "localhost",
        port: 8080,
        https:false,
        open: false,
        hotOnly:true,
        proxy:{
            "/apis" :{
                target:"http://127.0.0.1:10000",
                changeOrigin:true,
                ws:true,
                secure:false,
                pathRewrite:{
                    "^/apis":"/"
                }
            },
        },
    }
}