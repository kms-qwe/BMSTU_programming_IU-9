#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
import random
import subprocess
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path
from typing import Any


ROOT_DIR = Path(__file__).resolve().parent.parent
ENV_FILE = ROOT_DIR / "build" / ".env"
DEFAULT_API_BASE = "http://localhost:8080"


EMPLOYEES = [
    {"full_name": "Иван Петров", "phone": "+79990000001", "login": "ivan", "password": "123456"},
    {"full_name": "Анна Смирнова", "phone": "+79990000002", "login": "anna", "password": "123456"},
    {"full_name": "Петр Соколов", "phone": "+79990000003", "login": "petr", "password": "123456"},
    {"full_name": "Мария Волкова", "phone": "+79990000004", "login": "maria", "password": "123456"},
    {"full_name": "Алексей Морозов", "phone": "+79990000005", "login": "alex", "password": "123456"},
    {"full_name": "Ольга Федорова", "phone": "+79990000006", "login": "olga", "password": "123456"},
]

INGREDIENTS = [
    {"name": "Молоко", "unit": "л", "initial_stock": 800},
    {"name": "Кофейные зерна", "unit": "г", "initial_stock": 80000},
    {"name": "Вода", "unit": "мл", "initial_stock": 200000},
    {"name": "Сироп ванильный", "unit": "мл", "initial_stock": 30000},
    {"name": "Сироп карамельный", "unit": "мл", "initial_stock": 30000},
    {"name": "Сироп шоколадный", "unit": "мл", "initial_stock": 30000},
    {"name": "Сливки", "unit": "л", "initial_stock": 200},
    {"name": "Сахар", "unit": "г", "initial_stock": 50000},
    {"name": "Лед", "unit": "г", "initial_stock": 100000},
    {"name": "Матча", "unit": "г", "initial_stock": 10000},
    {"name": "Чай черный", "unit": "г", "initial_stock": 5000},
    {"name": "Основа для лимонада", "unit": "мл", "initial_stock": 50000},
    {"name": "Апельсиновый сок", "unit": "мл", "initial_stock": 50000},
    {"name": "Какао", "unit": "г", "initial_stock": 15000},
    {"name": "Специи для фильтра", "unit": "г", "initial_stock": 1200},
    {"name": "Топпинг фисташковый", "unit": "мл", "initial_stock": 0},
]

MENU_ITEMS = [
    {"name": "Эспрессо", "price": 170, "recipe": [("Кофейные зерна", 18), ("Вода", 60)]},
    {"name": "Американо", "price": 190, "recipe": [("Кофейные зерна", 18), ("Вода", 180)]},
    {"name": "Капучино", "price": 240, "recipe": [("Кофейные зерна", 18), ("Молоко", 0.18)]},
    {"name": "Латте", "price": 250, "recipe": [("Кофейные зерна", 18), ("Молоко", 0.25)]},
    {"name": "Флэт Уайт", "price": 260, "recipe": [("Кофейные зерна", 20), ("Молоко", 0.20)]},
    {"name": "Ванильный Латте", "price": 290, "recipe": [("Кофейные зерна", 18), ("Молоко", 0.25), ("Сироп ванильный", 25)]},
    {"name": "Карамельный Раф", "price": 320, "recipe": [("Кофейные зерна", 18), ("Сливки", 0.18), ("Сироп карамельный", 25), ("Сахар", 10)]},
    {"name": "Мокка", "price": 310, "recipe": [("Кофейные зерна", 18), ("Молоко", 0.22), ("Сироп шоколадный", 30)]},
    {"name": "Айс Латте", "price": 280, "recipe": [("Кофейные зерна", 18), ("Молоко", 0.22), ("Лед", 120)]},
    {"name": "Матча Латте", "price": 300, "recipe": [("Матча", 4), ("Молоко", 0.25), ("Сироп ванильный", 15)]},
    {"name": "Черный Чай", "price": 160, "recipe": [("Чай черный", 4), ("Вода", 300), ("Сахар", 5)]},
    {"name": "Цитрусовый Лимонад", "price": 270, "recipe": [("Основа для лимонада", 120), ("Апельсиновый сок", 80), ("Лед", 150)]},
    {"name": "Какао", "price": 230, "recipe": [("Какао", 25), ("Молоко", 0.25), ("Сахар", 10)]},
    {"name": "Фильтр Кофе", "price": 210, "recipe": [("Кофейные зерна", 24), ("Вода", 250), ("Специи для фильтра", 1)]},
]

EXTRA_MENU_ITEM = {
    "name": "Сезонный Латте",
    "price": 340,
    "recipe": [("Кофейные зерна", 18), ("Молоко", 0.25), ("Сироп карамельный", 30)],
}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Seed demo data for coffee shop backend")
    parser.add_argument("--orders", type=int, default=1000, help="number of orders to create")
    parser.add_argument("--api-base", default=DEFAULT_API_BASE, help="backend base url")
    return parser.parse_args()


def load_env() -> dict[str, str]:
    if not ENV_FILE.exists():
        raise SystemExit(f"env file not found: {ENV_FILE}")

    env: dict[str, str] = {}
    for line in ENV_FILE.read_text(encoding="utf-8").splitlines():
        stripped = line.strip()
        if not stripped or stripped.startswith("#") or "=" not in stripped:
            continue
        key, value = stripped.split("=", 1)
        env[key.strip()] = value.strip()
    return env


def run_command(args: list[str], *, input_text: str | None = None) -> str:
    result = subprocess.run(
        args,
        cwd=ROOT_DIR,
        input=input_text,
        capture_output=True,
        text=True,
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(f"command failed: {' '.join(args)}\nstdout:\n{result.stdout}\nstderr:\n{result.stderr}")
    return result.stdout.strip()


def docker_compose_cmd(env_file: Path) -> list[str]:
    return ["docker", "compose", "--env-file", str(env_file), "-f", str(ROOT_DIR / "build" / "docker-compose.yml")]


def cleanup_database(env: dict[str, str]) -> None:
    print("Cleaning database...")
    sql = """
TRUNCATE TABLE
  inventory_operations,
  order_items,
  orders,
  recipe_ingredients,
  ingredients,
  menu_items,
  shifts,
  employees
RESTART IDENTITY CASCADE;
"""
    run_command(
        docker_compose_cmd(ENV_FILE)
        + [
            "exec",
            "-T",
            "postgres",
            "psql",
            "-v",
            "ON_ERROR_STOP=1",
            "-U",
            env["DB_USER"],
            "-d",
            env["DB_NAME"],
            "-c",
            sql,
        ]
    )


def adjust_historical_timestamps(env: dict[str, str]) -> None:
    print("Adjusting timestamps for more realistic historical data...")
    sql = """
WITH closed_shifts AS (
  SELECT id, row_number() OVER (ORDER BY opened_at, id) AS rn
  FROM shifts
  WHERE is_active = false
),
scheduled AS (
  SELECT
    id,
    now() - make_interval(days => ((rn * 4)::int)) - interval '9 hours' AS new_opened_at,
    now() - make_interval(days => ((rn * 4)::int)) - interval '1 hour' AS new_closed_at
  FROM closed_shifts
)
UPDATE shifts s
SET opened_at = scheduled.new_opened_at,
    closed_at = scheduled.new_closed_at
FROM scheduled
WHERE s.id = scheduled.id;

WITH closed_shift_orders AS (
  SELECT
    o.id,
    s.opened_at,
    row_number() OVER (PARTITION BY o.shift_id ORDER BY o.id) AS rn,
    count(*) OVER (PARTITION BY o.shift_id) AS total
  FROM orders o
  JOIN shifts s ON s.id = o.shift_id
  WHERE s.is_active = false
)
UPDATE orders o
SET created_at = cso.opened_at + make_interval(mins => GREATEST(3, LEAST(420, (cso.rn * 420 / GREATEST(cso.total, 1))))::int)
FROM closed_shift_orders cso
WHERE o.id = cso.id;

UPDATE inventory_operations io
SET created_at = COALESCE(
  (SELECT o.created_at FROM orders o WHERE o.id = io.order_id),
  (SELECT s.opened_at + interval '30 minutes' FROM shifts s WHERE s.id = io.shift_id)
)
WHERE EXISTS (
  SELECT 1
  FROM shifts s
  WHERE s.id = io.shift_id
    AND s.is_active = false
);
"""
    run_command(
        docker_compose_cmd(ENV_FILE)
        + [
            "exec",
            "-T",
            "postgres",
            "psql",
            "-v",
            "ON_ERROR_STOP=1",
            "-U",
            env["DB_USER"],
            "-d",
            env["DB_NAME"],
            "-c",
            sql,
        ]
    )


def wait_for_api(base_url: str, timeout_seconds: int = 60) -> None:
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        try:
            request_json("GET", f"{base_url}/api/shift/active")
            return
        except Exception:
            time.sleep(1)
    raise RuntimeError(f"backend did not become ready within {timeout_seconds} seconds")


def request_json(method: str, url: str, payload: dict[str, Any] | None = None) -> Any:
    headers = {"Content-Type": "application/json"}
    body = None
    if payload is not None:
        body = json.dumps(payload).encode("utf-8")

    request = urllib.request.Request(url, data=body, headers=headers, method=method)
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            raw = response.read()
            if not raw:
                return None
            return json.loads(raw.decode("utf-8"))
    except urllib.error.HTTPError as exc:
        try:
            details = exc.read().decode("utf-8")
        except Exception:
            details = "<failed to read body>"
        raise RuntimeError(f"{method} {url} failed with status {exc.code}: {details}") from exc


def create_employees() -> None:
    print("Creating employees...")
    script_path = ROOT_DIR / "scripts" / "create_employee.sh"
    for employee in EMPLOYEES:
        run_command(
            [
                str(script_path),
                employee["full_name"],
                employee["phone"],
                employee["login"],
                employee["password"],
            ]
        )


def open_shift(base_url: str, login: str, password: str) -> None:
    request_json("POST", f"{base_url}/api/shift/open", {"login": login, "password": password})


def close_shift(base_url: str) -> None:
    request_json("POST", f"{base_url}/api/shift/close", {})


def create_ingredients(base_url: str) -> dict[str, str]:
    print("Creating ingredients...")
    ingredient_ids: dict[str, str] = {}
    for ingredient in INGREDIENTS:
        response = request_json(
            "POST",
            f"{base_url}/api/ingredient",
            {
                "name": ingredient["name"],
                "unit": ingredient["unit"],
                "initial_stock": ingredient["initial_stock"],
            },
        )
        ingredient_ids[ingredient["name"]] = response["ingredient"]["id"]
    return ingredient_ids


def create_menu_items(base_url: str, ingredient_ids: dict[str, str]) -> tuple[dict[str, str], str]:
    print("Creating menu items...")
    menu_item_ids: dict[str, str] = {}
    for menu_item in MENU_ITEMS:
        response = request_json(
            "POST",
            f"{base_url}/api/menu-item",
            {
                "name": menu_item["name"],
                "price": menu_item["price"],
                "recipe": [
                    {"ingredient_id": ingredient_ids[ingredient_name], "amount": amount}
                    for ingredient_name, amount in menu_item["recipe"]
                ],
            },
        )
        menu_item_ids[menu_item["name"]] = response["menu_item"]["id"]

    inactive_response = request_json(
        "POST",
        f"{base_url}/api/menu-item",
        {
            "name": EXTRA_MENU_ITEM["name"],
            "price": EXTRA_MENU_ITEM["price"],
            "recipe": [
                {"ingredient_id": ingredient_ids[ingredient_name], "amount": amount}
                for ingredient_name, amount in EXTRA_MENU_ITEM["recipe"]
            ],
        },
    )
    inactive_menu_item_id = inactive_response["menu_item"]["id"]
    return menu_item_ids, inactive_menu_item_id


def create_inventory_operation(base_url: str, ingredient_id: str, operation_type: str, amount: float) -> None:
    request_json(
        "POST",
        f"{base_url}/api/inventory/operation",
        {"ingredient_id": ingredient_id, "operation_type": operation_type, "amount": amount},
    )


def create_order(base_url: str, items: list[dict[str, Any]]) -> None:
    request_json("POST", f"{base_url}/api/order", {"items": items})


def deactivate_ingredient(base_url: str, ingredient_id: str) -> None:
    request_json("PATCH", f"{base_url}/api/ingredient/{ingredient_id}/activity", {"is_active": False})


def deactivate_menu_item(base_url: str, menu_item_id: str) -> None:
    request_json("PATCH", f"{base_url}/api/menu-item/{menu_item_id}/activity", {"is_active": False})


def build_shift_plan(orders_count: int) -> list[tuple[str, str, int]]:
    weights = [18, 17, 16, 17, 16, 16]
    counts = [(orders_count * weight) // 100 for weight in weights]
    remainder = orders_count - sum(counts)
    for index in range(remainder):
        counts[index % len(counts)] += 1

    return [
        (employee["login"], employee["password"], counts[index])
        for index, employee in enumerate(EMPLOYEES)
    ]


def seed_orders_and_operations(base_url: str, orders_count: int, ingredient_ids: dict[str, str], menu_item_ids: dict[str, str]) -> str:
    rng = random.Random(42)
    shift_plan = build_shift_plan(orders_count)

    weighted_menu = [
        ("Эспрессо", 8),
        ("Американо", 12),
        ("Капучино", 16),
        ("Латте", 17),
        ("Флэт Уайт", 10),
        ("Ванильный Латте", 9),
        ("Карамельный Раф", 8),
        ("Мокка", 7),
        ("Айс Латте", 8),
        ("Матча Латте", 6),
        ("Черный Чай", 5),
        ("Цитрусовый Лимонад", 6),
        ("Какао", 5),
        ("Фильтр Кофе", 7),
    ]
    menu_names = [name for name, _ in weighted_menu]
    menu_weights = [weight for _, weight in weighted_menu]

    created_orders = 0
    active_login = EMPLOYEES[0]["login"]
    for shift_index, (login, password, shift_orders) in enumerate(shift_plan, start=1):
        if shift_orders <= 0:
            continue

        if shift_index == 1:
            print(f"Using bootstrap shift for {login} with {shift_orders} orders...")
        else:
            print(f"Opening shift {shift_index} for {login} with {shift_orders} orders...")
            open_shift(base_url, login, password)
            active_login = login

        create_inventory_operation(base_url, ingredient_ids["Молоко"], "MANUAL_RESTOCK", 25)
        create_inventory_operation(base_url, ingredient_ids["Кофейные зерна"], "MANUAL_RESTOCK", 5000)

        for order_index in range(shift_orders):
            items_count = rng.choices([1, 2, 3, 4], weights=[35, 40, 20, 5], k=1)[0]
            chosen_names = rng.choices(menu_names, weights=menu_weights, k=items_count)
            aggregated: dict[str, int] = {}
            for chosen_name in chosen_names:
                aggregated[chosen_name] = aggregated.get(chosen_name, 0) + rng.choices([1, 2, 3], weights=[70, 25, 5], k=1)[0]

            payload_items = [
                {"menu_item_id": menu_item_ids[name], "quantity": quantity}
                for name, quantity in aggregated.items()
            ]
            create_order(base_url, payload_items)
            created_orders += 1

            if order_index > 0 and order_index % 40 == 0:
                create_inventory_operation(base_url, ingredient_ids["Молоко"], "MANUAL_RESTOCK", rng.choice([10, 15, 20]))
            if order_index > 0 and order_index % 55 == 0:
                create_inventory_operation(base_url, ingredient_ids["Сироп ванильный"], "MANUAL_RESTOCK", rng.choice([500, 750, 1000]))
            if order_index > 0 and order_index % 70 == 0:
                create_inventory_operation(base_url, ingredient_ids["Сахар"], "MANUAL_WRITE_OFF", rng.choice([100, 150, 200]))

            if created_orders % 100 == 0:
                print(f"Created {created_orders} orders...")

        has_next_non_empty_shift = any(future_orders > 0 for _, _, future_orders in shift_plan[shift_index:])
        if has_next_non_empty_shift:
            close_shift(base_url)

    return active_login


def main() -> None:
    args = parse_args()
    if args.orders < 1:
        raise SystemExit("--orders must be greater than zero")

    env = load_env()
    wait_for_api(args.api_base)
    cleanup_database(env)
    create_employees()

    print("Opening bootstrap shift...")
    open_shift(args.api_base, EMPLOYEES[0]["login"], EMPLOYEES[0]["password"])

    ingredient_ids = create_ingredients(args.api_base)
    menu_item_ids, inactive_menu_item_id = create_menu_items(args.api_base, ingredient_ids)

    active_login = seed_orders_and_operations(args.api_base, args.orders, ingredient_ids, menu_item_ids)

    deactivate_ingredient(args.api_base, ingredient_ids["Топпинг фисташковый"])
    deactivate_menu_item(args.api_base, inactive_menu_item_id)
    adjust_historical_timestamps(env)

    print()
    print("Demo data is ready.")
    print(f"API base: {args.api_base}")
    print(f"Orders created: {args.orders}")
    print("Employees:")
    for employee in EMPLOYEES:
        print(f"  - {employee['login']} / {employee['password']} ({employee['full_name']})")
    print(f"Active shift remains open for: {active_login} / 123456")


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit("Interrupted")
