"""Scaffold a new microservice from the hello template.

Usage:
    cd server/
    uv run python scripts/scaffold.py <service_name>

Example:
    uv run python scripts/scaffold.py auth
"""

import os
import re
import shutil
import sys

SERVICES_DIR = os.path.join(os.path.dirname(__file__), "..", "services")
HELLO_DIR = os.path.join(SERVICES_DIR, "hello")


def validate_name(name: str) -> str:
    name = name.strip().lower()
    if not re.match(r"^[a-z][a-z0-9_]*$", name):
        print(
            f"Error: '{name}' is not a valid service name. Use lowercase letters, digits, underscores, starting with a letter."
        )
        sys.exit(1)
    if os.path.exists(os.path.join(SERVICES_DIR, name)):
        print(f"Error: services/{name}/ already exists.")
        sys.exit(1)
    return name


def copy_template(name: str) -> str:
    target = os.path.join(SERVICES_DIR, name)
    shutil.copytree(HELLO_DIR, target, symlinks=False)
    return target


def rename_cmd_dir(target: str, name: str) -> None:
    old_cmd = os.path.join(target, "cmd", "hello")
    new_cmd = os.path.join(target, "cmd", name)
    if os.path.exists(old_cmd):
        os.renames(old_cmd, new_cmd)


def replace_in_file(filepath: str, old: str, new: str) -> None:
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(content)


def walk_and_replace(target: str, name: str) -> None:
    old_module = "github.com/justblue/luoye/services/hello"
    new_module = f"github.com/justblue/luoye/services/{name}"
    old_app_name = 'kratos.Name("hello")'
    new_app_name = f'kratos.Name("{name}")'

    for root, dirs, files in os.walk(target):
        for fname in files:
            fpath = os.path.join(root, fname)
            if fname == "go.sum":
                os.remove(fpath)
                continue
            replace_in_file(fpath, old_module, new_module)
            replace_in_file(fpath, old_app_name, new_app_name)


def update_go_mod(target: str, name: str) -> None:
    gomod = os.path.join(target, "go.mod")
    if os.path.exists(gomod):
        with open(gomod, "r", encoding="utf-8") as f:
            content = f.read()
        content = content.replace(
            "module github.com/justblue/luoye/services/hello",
            f"module github.com/justblue/luoye/services/{name}",
        )
        with open(gomod, "w", encoding="utf-8") as f:
            f.write(content)


def main():
    if len(sys.argv) != 2:
        print(__doc__.strip())
        sys.exit(1)

    name = validate_name(sys.argv[1])
    target = copy_template(name)
    rename_cmd_dir(target, name)
    walk_and_replace(target, name)
    update_go_mod(target, name)

    print(f"Created services/{name}/ from template.")
    print()
    print("Next steps:")
    print(f"  1. cd server")
    print(f"  2. go work use ./services/{name}")
    print(
        f"  3. go mod edit -module github.com/justblue/luoye/services/{name} ./services/{name}/go.mod  # if needed"
    )
    print(f"  4. cd services/{name} && go mod tidy")
    print(f"  5. Edit internal/domain/*.go with your business logic")
    print(f"  6. Update gRPC port in config/config.yaml (avoid conflicts)")
    print(f"  7. Add gRPC proxy in gateway/internal/proxy/{name}.go")
    print(f"  8. Register routes in gateway/internal/server/http.go")


if __name__ == "__main__":
    main()
