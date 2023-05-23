module.exports = {
  apps : [
    {
      name: "Proxy",
      script: "proxy.js",
      watch: true,
    },
    {
      name: "Frontend",
      script: "server.js",
      watch: true,
    }
  ]
}