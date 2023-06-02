// get url by env, if dev, return xxx:8081, if prod, return xxx:8080
const getEnvUrl = () => {
    if (process.env.NODE_ENV === 'development') {
        return 'http://localhost:8080'
    } else {
        return window.location.origin
    }
}

export {
    getEnvUrl
}