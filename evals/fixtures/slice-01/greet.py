def greet(name: str) -> str:
    if not name.strip():
        raise ValueError("name must not be blank")
    return f"Hello, {name}"
