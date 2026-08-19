# How to Contribute

We're glad you'd like to help! Below is a brief guide to participating in the project.

## Before You Start

- Review the [README.md](README.md) to understand the project's goals and structure
- Check the [list of current issues and tasks](https://github.com/idfactory/idtools/issues) — your idea might already be under discussion
- Make sure you agree with the project's [license](LICENSE)

## How to Propose Changes

- Create a new issue [here](https://github.com/idfactory/idtools/issues) to discuss your proposed changes
- Mention the maintainer `@hex21h` in your issue, clearly describing what you want to change and why
- If you’ve agreed that the change is needed and you’ll implement it yourself, the process is straightforward:
  - Fork the repository (click the "Fork" button)
  - Create a dedicated branch for your work
  - Make your changes and test them thoroughly
  - Commit with a clear message, e.g.: `added feature X and improved Y`
  - Push your changes to your fork.
  - Open a pull request against the `main` branch, including:
    - A description of what was changed and why
    - Instructions on how to verify the changes
    - Links to any related issues
  - Wait for a code review and eventual merge
  - PROFIT!

## Code of Conduct and Style Guidelines

- Follow the project's established code and documentation style
- Be respectful everywhere: in comments, issues, or discussions. No insults, profanity, or inappropriate content
- If you believe the entire project needs a complete rewrite or a switch to different technologies, you’re probably mistaken
- Always be polite and respectful toward other contributors
- Just be a decent human being

## Got Questions?

If you have any questions, create an issue [here](https://github.com/idfactory/idtools/issues), mention `@hex21h`, and we’ll discuss it together.

## Tests

```sh
# run tests
go test -v ./...

# coverage
go test -cover ./...
```

## Releasing a New Version

```bash
git tag v1.0.0
git push origin main
git push origin v1.0.0
```

## Acknowledgments

We appreciate all contributions—whether it's code, documentation, testing, or ideas. Thank you!
