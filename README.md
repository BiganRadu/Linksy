# 🔗 Linksy

Linksy is a modern **URL shortener with analytics and QR code support**.  
It allows you to create, manage, and analyze your shortened links with an intuitive web interface.

🌍 **Live Demo**: [Linksy Project](http://linksyproject.s3-website.eu-north-1.amazonaws.com/)

---

## ✨ Features

- 📌 **Main Page** – project overview and information  
  _![img.png](img.png)_

- 🔗 **Links Page** – view, create, edit, and delete links  
  _![img_1.png](img_1.png)_

- 🖼 **QR Codes Page** – automatically generated QR codes for your links, with images stored in AWS S3  
  _![img_2.png](img_2.png)_

- 📊 **Analytics Page** – in-depth statistics:
    - Daily accesses chart
    - Round charts for **countries** and **platforms**
    - Detailed table with per-link stats  
      _![img_3.png](img_3.png)_
      _![img_4.png](img_4.png)_

- ⚙️ **Settings Page** – update your account details (username, password)  
  _![img_5.png](img_5.png)_

- 👤 **Authentication** – sign-up and sign-in pages for account management  
  _![img_6.png](img_6.png)_
  _![img_7.png](img_7.png)_

---

## 🛠️ Tech Stack

**Backend**
- Language: [Go](https://go.dev/)
- Build System: [Bazel](https://bazel.build/)
- Database: [MongoDB](https://www.mongodb.com/)
- Hosting: [Render](https://render.com/)

**Frontend**
- Framework: [React](https://reactjs.org/)
- UI Components: [MUI](https://mui.com/)
- Hosting: [AWS S3](https://aws.amazon.com/s3/) (static hosting)

**Storage**
- QR code images are generated in the backend and stored in **AWS S3**

