module.exports = {
  apps: [
    {
      name: "pms-api",
      script: "go",
      args: "run ./cmd/api/main.go",
      cwd: __dirname,
      max_restarts: 3,
      autorestart: true,
      watch: false,
    },
    {
      name: "pms-web",
      script: "npm",
      args: "run dev",
      cwd: "/Users/enaira/Desktop/ENTERPRISE PMS/pms-web",
      max_restarts: 3,
      autorestart: true,
      watch: false,
    },
  ],
};
