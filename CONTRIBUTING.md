### **Contributing to Beerus**

Thank you for your interest in contributing to **Beerus**! 🎉 We appreciate your help in making this project better. Please follow the guidelines below to ensure a smooth contribution process.

---

## Development Setup

### Prerequisites

- Go 1.23+
- Docker
- pre-commit
- golangci-lint

### Local Setup

1. **Fork the Repository**
   - Click the "Fork" button at the top right of this repository.
   - Clone your fork locally:
     ```sh
     git clone https://github.com/{your_username}/beerus.git
     cd beerus
     ```

2. **Create a New Branch**
   - Use a descriptive name for your branch:
     ```sh
     git checkout -b feature/your-feature-name
     ```

3. **Install Dependencies**
   - Ensure you have the necessary dependencies installed:
     ```sh
        go mod download
     ```

4. **Make Your Changes**
   - Follow the project's coding style and best practices.
   - Write clear, maintainable, and well-documented code.

5. **Run Tests**
   - Before submitting, make sure your changes don’t break existing functionality:
     ```sh
       go test ./...
     ```

6. **Commit Your Changes**
   - Use meaningful commit messages:
     ```sh
     git commit -m "feat: add new feature description"
     ```
   - Push your changes to your forked repository:
     ```sh
     git push origin feature/your-feature-name
     ```

7. **Create a Pull Request**
   - Open a PR against the `main` branch of this repository.
   - Provide a **clear description** of the changes you made.
   - If applicable, include screenshots or logs to help reviewers.

---

## **Code Guidelines**

- Follow the existing **coding style** and project conventions.
- Write **unit tests** for new features and bug fixes.
- Keep pull requests **small and focused**.
- Document any new public functions or APIs.

---

## **Reporting Issues**

If you find a bug or have a feature request, please **open an issue**:

1. Check if the issue already exists in [existing issues](https://github.com/lucasmendesl/beerus/issues).
2. If not, create a new issue with a clear **title** and **description**.
3. Provide steps to reproduce (if it's a bug) or a clear use case (if it's a feature request).

---

## **Code of Conduct**

We follow the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/). Please be respectful and inclusive to maintain a welcoming community.

---

🚀 **Thank you for contributing to Beerus!** We appreciate your help in making this project awesome! 💙
