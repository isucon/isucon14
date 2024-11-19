import nox


@nox.session(python="3.10")
def lint(session: nox.Session) -> None:
    session.install("pre-commit")
    session.run("pre-commit", "run", "--all-files")


@nox.session(python="3.10")
def mypy(session: nox.Session) -> None:
    session.install(".")
    session.install("mypy")
    session.run(
        "mypy",
        "-p",
        "isuride",
    )
