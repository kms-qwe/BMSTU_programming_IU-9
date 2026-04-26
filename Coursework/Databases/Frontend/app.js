const API_BASE_URL = "";

const ROUTES = {
  shift: "/",
  orders: "/orders",
  menu: "/menu",
  ingredients: "/ingredients",
  inventory: "/inventory",
  reports: "/reports",
};

const MANUAL_OPERATION_TYPES = [
  { value: "MANUAL_RESTOCK", label: "Пополнение" },
  { value: "MANUAL_WRITE_OFF", label: "Списание" },
];

const OPERATION_TYPE_OPTIONS = [
  { value: "AUTO_ORDER_WRITE_OFF", label: "Автоматическое списание по заказу" },
  ...MANUAL_OPERATION_TYPES,
  { value: "INITIAL_STOCK", label: "Начальный остаток" },
];

const REPORT_FILES = {
  sales: "sales-report.pdf",
  employees: "employees-report.pdf",
  inventory: "inventory-report.pdf",
  popularity: "menu-popularity-report.pdf",
};

const state = {
  activeShift: null,
};

const appRoot = document.querySelector("#app");
const modalRoot = document.querySelector("#modal-root");
const toastRoot = document.querySelector("#toast-root");

class ApiError extends Error {
  constructor(payload) {
    super(payload.message || "Ошибка запроса");
    this.code = payload.code || "UNKNOWN_ERROR";
    this.details = payload.details ?? null;
  }
}

const escapeHtml = (value) =>
  String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");

const formatDateTime = (value) => {
  if (!value) {
    return "—";
  }
  return new Intl.DateTimeFormat("ru-RU", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(new Date(value));
};

const toDateTimeLocal = (value) => {
  if (!value) {
    return "";
  }
  const date = new Date(value);
  const offset = date.getTimezoneOffset();
  const normalized = new Date(date.getTime() - offset * 60000);
  return normalized.toISOString().slice(0, 16);
};

const formatMoney = (value) =>
  new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    minimumFractionDigits: 2,
  }).format(Number(value ?? 0));

const formatNumber = (value) =>
  new Intl.NumberFormat("ru-RU", {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  }).format(Number(value ?? 0));

const badge = (isActive) =>
  `<span class="badge ${isActive ? "active" : "inactive"}">${isActive ? "Активен" : "Неактивен"}</span>`;

const formatOperationType = (value) =>
  OPERATION_TYPE_OPTIONS.find((item) => item.value === value)?.label || value || "—";

const showToast = (message, type = "error") => {
  const toast = document.createElement("div");
  toast.className = `toast ${type}`;
  toast.textContent = message;
  toastRoot.appendChild(toast);
  setTimeout(() => toast.remove(), 3600);
};

const getErrorMessage = (error) => {
  if (!(error instanceof ApiError)) {
    return "Не удалось выполнить запрос";
  }
  const messages = {
    DUPLICATE_NAME: "Сущность с таким названием уже существует",
    INACTIVE_INGREDIENT: "Нельзя использовать неактивный ингредиент в рецепте",
    INACTIVE_MENU_ITEM: "Нельзя добавить неактивный пункт меню",
    INSUFFICIENT_STOCK: "Недостаточно ингредиентов",
    UNSUPPORTED_OPERATION: "Неподдерживаемый тип операции",
    INVALID_CREDENTIALS: "Неверный логин или пароль",
    NOT_FOUND: "Сущность не найдена",
    INTERNAL_ERROR: "Внутренняя ошибка сервера",
  };

  if (error.code === "VALIDATION_ERROR") {
    return error.message || "Ошибка валидации";
  }

  return messages[error.code] || error.message || "Произошла ошибка";
};

const handleApiError = (error) => {
  if (error instanceof ApiError && error.code === "SHIFT_NOT_ACTIVE") {
    state.activeShift = null;
    if (window.location.pathname !== ROUTES.shift) {
      navigate(ROUTES.shift, { replace: true });
    } else {
      renderApp();
    }
    return;
  }
  showToast(getErrorMessage(error), "error");
};

const buildQuery = (query) => {
  const params = new URLSearchParams();
  Object.entries(query || {}).forEach(([key, value]) => {
    if (value !== "" && value !== undefined && value !== null) {
      params.set(key, value);
    }
  });
  const queryString = params.toString();
  return queryString ? `?${queryString}` : "";
};

const apiRequest = async (path, options = {}) => {
  const { method = "GET", body, query, responseType = "json" } = options;
  const response = await fetch(`${API_BASE_URL}${path}${buildQuery(query)}`, {
    method,
    headers: {
      ...(responseType === "json" ? { Accept: "application/json" } : {}),
      ...(body ? { "Content-Type": "application/json" } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!response.ok) {
    let payload = null;
    const contentType = response.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      payload = await response.json();
    }
    const errorBody = payload?.error || {};
    throw new ApiError({
      code: errorBody.code,
      message: errorBody.message,
      details: errorBody.details,
    });
  }

  if (responseType === "blob") {
    return response.blob();
  }

  if (response.status === 204) {
    return null;
  }

  const contentType = response.headers.get("content-type") || "";
  if (!contentType.includes("application/json")) {
    return null;
  }
  return response.json();
};

const api = {
  getActiveShift: () => apiRequest("/api/shift/active"),
  openShift: (payload) => apiRequest("/api/shift/open", { method: "POST", body: payload }),
  closeShift: () => apiRequest("/api/shift/close", { method: "POST" }),
  getEmployees: () => apiRequest("/api/employee"),
  getIngredients: (includeInactive = true) =>
    apiRequest("/api/ingredient", { query: { include_inactive: includeInactive } }),
  getIngredient: (ingredientId) => apiRequest(`/api/ingredient/${ingredientId}`),
  createIngredient: (payload) => apiRequest("/api/ingredient", { method: "POST", body: payload }),
  setIngredientActivity: (ingredientId, isActive) =>
    apiRequest(`/api/ingredient/${ingredientId}/activity`, {
      method: "PATCH",
      body: { is_active: isActive },
    }),
  getMenuItems: (includeInactive = true) =>
    apiRequest("/api/menu-item", { query: { include_inactive: includeInactive } }),
  getMenuItem: (menuItemId) => apiRequest(`/api/menu-item/${menuItemId}`),
  createMenuItem: (payload) => apiRequest("/api/menu-item", { method: "POST", body: payload }),
  updateMenuPrice: (menuItemId, price) =>
    apiRequest(`/api/menu-item/${menuItemId}/price`, { method: "PATCH", body: { price } }),
  setMenuActivity: (menuItemId, isActive) =>
    apiRequest(`/api/menu-item/${menuItemId}/activity`, {
      method: "PATCH",
      body: { is_active: isActive },
    }),
  getOrders: (query) => apiRequest("/api/order", { query }),
  getOrder: (orderId) => apiRequest(`/api/order/${orderId}`),
  createOrder: (payload) => apiRequest("/api/order", { method: "POST", body: payload }),
  getOperations: (query) => apiRequest("/api/inventory/operation", { query }),
  createOperation: (payload) =>
    apiRequest("/api/inventory/operation", { method: "POST", body: payload }),
  downloadSalesReport: (payload) =>
    apiRequest("/api/report/sales/pdf", { method: "POST", body: payload, responseType: "blob" }),
  downloadEmployeesReport: (payload) =>
    apiRequest("/api/report/employees/pdf", {
      method: "POST",
      body: payload,
      responseType: "blob",
    }),
  downloadInventoryReport: (payload) =>
    apiRequest("/api/report/inventory/pdf", {
      method: "POST",
      body: payload,
      responseType: "blob",
    }),
  downloadPopularityReport: (payload) =>
    apiRequest("/api/report/menu-popularity/pdf", {
      method: "POST",
      body: payload,
      responseType: "blob",
    }),
};

const navigate = (path, { replace = false } = {}) => {
  if (replace) {
    window.history.replaceState({}, "", path);
  } else {
    window.history.pushState({}, "", path);
  }
  renderApp();
};

const openModal = (content, { small = false } = {}) => {
  modalRoot.innerHTML = `
    <div class="modal-backdrop">
      <div class="modal ${small ? "small" : ""}">
        ${content}
      </div>
    </div>
  `;
  modalRoot.querySelector(".modal-backdrop").addEventListener("click", (event) => {
    if (event.target.classList.contains("modal-backdrop")) {
      closeModal();
    }
  });
};

const closeModal = () => {
  modalRoot.innerHTML = "";
};

const setApp = (content) => {
  appRoot.innerHTML = content;
};

const renderLoading = (message = "Загрузка...") => {
  setApp(`<div class="auth-page"><div class="status">${message}</div></div>`);
};

const renderOpenShiftScreen = () => {
  setApp(`
    <div class="auth-page">
      <div class="auth-card">
        <div class="title-row">
          <div>
            <h1>Открытие смены</h1>
            <p class="muted">Начните смену, чтобы перейти к работе с заказами, меню и складом.</p>
          </div>
        </div>
        <form id="open-shift-form" class="grid">
          <div class="field">
            <label for="login">Логин</label>
            <input id="login" name="login" required />
          </div>
          <div class="field">
            <label for="password">Пароль</label>
            <input id="password" name="password" type="password" required />
          </div>
          <button class="btn primary" type="submit">Открыть смену</button>
        </form>
      </div>
    </div>
  `);

  document.querySelector("#open-shift-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const submitButton = event.currentTarget.querySelector("button[type='submit']");
    submitButton.disabled = true;
    const formData = new FormData(event.currentTarget);

    try {
      const response = await api.openShift({
        login: formData.get("login"),
        password: formData.get("password"),
      });
      state.activeShift = response.shift;
      navigate(ROUTES.orders, { replace: true });
    } catch (error) {
      if (error instanceof ApiError && error.code === "INVALID_CREDENTIALS") {
        showToast("Неверный логин или пароль", "error");
      } else if (error instanceof ApiError && error.code === "SHIFT_ALREADY_OPENED") {
        const active = await api.getActiveShift();
        if (active.is_active) {
          state.activeShift = active.shift;
          navigate(ROUTES.orders, { replace: true });
        }
      } else {
        handleApiError(error);
      }
    } finally {
      submitButton.disabled = false;
    }
  });
};

const routeTitle = (path) =>
  ({
    [ROUTES.orders]: "Заказы",
    [ROUTES.menu]: "Меню",
    [ROUTES.ingredients]: "Ингредиенты",
    [ROUTES.inventory]: "Складские операции",
    [ROUTES.reports]: "Отчеты",
  })[path] || "Coffee Shop";

const renderProtectedLayout = () => {
  const currentPath = window.location.pathname;
  const shift = state.activeShift;
  setApp(`
    <div class="page-shell">
      <aside class="sidebar">
        <div class="brand">
          <h1>Coffee Shop</h1>
          <p>Заказы, меню, склад и отчеты в одном месте</p>
        </div>
        <nav class="nav">
          ${Object.entries({
            [ROUTES.orders]: "Заказы",
            [ROUTES.menu]: "Меню",
            [ROUTES.ingredients]: "Ингредиенты",
            [ROUTES.inventory]: "Складские операции",
            [ROUTES.reports]: "Отчеты",
          })
            .map(
              ([path, label]) =>
                `<a class="nav-link ${currentPath === path ? "active" : ""}" href="${path}" data-nav>${label}</a>`,
            )
            .join("")}
        </nav>
        <div class="shift-meta">
          <div><strong>Смена активна</strong></div>
          <div>Сотрудник: ${escapeHtml(shift?.duty?.full_name || "—")}</div>
          <div>Открыта: ${formatDateTime(shift?.opened_at)}</div>
        </div>
      </aside>
      <main class="main">
        <div class="header">
          <div class="header-card">
            <h2>${routeTitle(currentPath)}</h2>
            <div class="shift-meta">Активная смена открыта ${formatDateTime(shift?.opened_at)}</div>
          </div>
          <div class="header-actions">
            <div class="header-card">
              <div><strong>${escapeHtml(shift?.duty?.full_name || "—")}</strong></div>
              <div class="shift-meta">${escapeHtml(shift?.duty?.login || "")}</div>
            </div>
            <button id="close-shift-button" class="btn danger">Закрыть смену</button>
          </div>
        </div>
        <div id="page-content"></div>
      </main>
    </div>
  `);

  document.querySelectorAll("[data-nav]").forEach((link) => {
    link.addEventListener("click", (event) => {
      event.preventDefault();
      navigate(link.getAttribute("href"));
    });
  });

  document.querySelector("#close-shift-button").addEventListener("click", async () => {
    const button = document.querySelector("#close-shift-button");
    button.disabled = true;
    try {
      await api.closeShift();
      state.activeShift = null;
      navigate(ROUTES.shift, { replace: true });
    } catch (error) {
      handleApiError(error);
      button.disabled = false;
    }
  });
};

const renderPageContent = (content) => {
  const container = document.querySelector("#page-content");
  if (container) {
    container.innerHTML = content;
  }
};

const renderStatusBlock = (selector, message) => {
  const element = document.querySelector(selector);
  if (element) {
    element.innerHTML = `<div class="status">${message}</div>`;
  }
};

const createEmptyState = (title, description) => `
  <div class="empty-state">
    <h3>${title}</h3>
    <p class="muted">${description}</p>
  </div>
`;

const downloadBlob = (blob, fileName) => {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = fileName;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
};

const getDefaultDateRange = () => {
  const now = new Date();
  const past = new Date(now.getTime() - 24 * 60 * 60 * 1000);
  return {
    from: toDateTimeLocal(past),
    to: toDateTimeLocal(now),
  };
};

const buildPaginationInfo = (pagination) => {
  if (!pagination) {
    return "";
  }
  return `
    <div class="pagination">
      Страница ${pagination.page} из ${pagination.total_pages}. Всего записей: ${pagination.total_items}, размер страницы: ${pagination.page_size}.
    </div>
  `;
};

const buildPaginationControls = (pagination, prefix) => {
  if (!pagination || pagination.total_pages <= 1) {
    return "";
  }

  const currentPage = Number(pagination.page || 1);
  const totalPages = Number(pagination.total_pages || 1);
  const pages = [];
  const start = Math.max(1, currentPage - 2);
  const end = Math.min(totalPages, currentPage + 2);

  for (let page = start; page <= end; page += 1) {
    pages.push(`
      <button
        type="button"
        class="btn ${page === currentPage ? "primary" : ""}"
        data-pagination-prefix="${prefix}"
        data-page="${page}"
        ${page === currentPage ? "disabled" : ""}
      >
        ${page}
      </button>
    `);
  }

  return `
    <div class="btn-row pagination-controls">
      <button
        type="button"
        class="btn"
        data-pagination-prefix="${prefix}"
        data-page="${currentPage - 1}"
        ${currentPage <= 1 ? "disabled" : ""}
      >
        Назад
      </button>
      ${pages.join("")}
      <button
        type="button"
        class="btn"
        data-pagination-prefix="${prefix}"
        data-page="${currentPage + 1}"
        ${currentPage >= totalPages ? "disabled" : ""}
      >
        Вперед
      </button>
    </div>
  `;
};

const ensureActiveShift = async () => {
  const active = await api.getActiveShift();
  if (!active.is_active) {
    state.activeShift = null;
    navigate(ROUTES.shift, { replace: true });
    return false;
  }
  state.activeShift = active.shift;
  return true;
};

const renderOrdersPage = async () => {
  let currentPage = 1;

  renderPageContent(`
    <div class="section-head">
      <div>
        <h3>Список заказов</h3>
        <p class="muted">Создание и просмотр без редактирования и удаления.</p>
      </div>
      <button class="btn primary" id="create-order-button">Создать заказ</button>
    </div>
    <form id="orders-filter-form" class="filter-bar">
      <div class="form-grid">
        <div class="field">
          <label for="orders-employee-filter">Сотрудник</label>
          <select id="orders-employee-filter" name="employee_id"></select>
        </div>
        <div class="field">
          <label for="orders-from">Дата/время from</label>
          <input id="orders-from" name="from" type="datetime-local" />
        </div>
        <div class="field">
          <label for="orders-to">Дата/время to</label>
          <input id="orders-to" name="to" type="datetime-local" />
        </div>
      </div>
      <div class="btn-row">
        <button type="button" class="btn" id="orders-filters-reset">Очистить фильтры</button>
      </div>
    </form>
    <div id="orders-table-block" class="table-wrap"><div class="status">Загрузка заказов...</div></div>
  `);

  const filterForm = document.querySelector("#orders-filter-form");
  const employeeSelect = document.querySelector("#orders-employee-filter");

  const loadOrders = async () => {
    renderStatusBlock("#orders-table-block", "Загрузка заказов...");
    try {
      const query = {
        page: currentPage,
        employee_id: filterForm.elements.employee_id.value,
        from: filterForm.elements.from.value ? new Date(filterForm.elements.from.value).toISOString() : "",
        to: filterForm.elements.to.value ? new Date(filterForm.elements.to.value).toISOString() : "",
      };
      const response = await api.getOrders(query);
      const rows = response.items || [];
      const content = rows.length
        ? `
            <table>
              <thead>
                <tr>
                  <th>Дата создания</th>
                  <th>Сотрудник</th>
                  <th>Позиции</th>
                  <th>Итоговая сумма</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                ${rows
                  .map(
                    (order) => `
                      <tr>
                        <td>${formatDateTime(order.created_at)}</td>
                        <td>${escapeHtml(order.employee?.full_name || "—")}</td>
                        <td>${escapeHtml(
                          (order.items || [])
                            .map((item) => `${item.menu_item_name} x${item.quantity}`)
                            .join(", "),
                        )}</td>
                        <td>${formatMoney(order.total_price)}</td>
                        <td><button class="btn" data-order-details="${order.id}">Подробнее</button></td>
                      </tr>
                    `,
                  )
                  .join("")}
              </tbody>
            </table>
            ${buildPaginationInfo(response.pagination)}
            ${buildPaginationControls(response.pagination, "orders")}
          `
        : createEmptyState("Заказов пока нет", "Создайте первый заказ для активной смены.");
      document.querySelector("#orders-table-block").innerHTML = content;
      document.querySelectorAll("[data-order-details]").forEach((button) => {
        button.addEventListener("click", () => openOrderDetails(button.dataset.orderDetails));
      });
      document.querySelectorAll('[data-pagination-prefix="orders"]').forEach((button) => {
        button.addEventListener("click", async () => {
          currentPage = Number(button.dataset.page);
          await loadOrders();
        });
      });
    } catch (error) {
      handleApiError(error);
      renderStatusBlock("#orders-table-block", "Не удалось загрузить заказы");
    }
  };

  const [employeesResponse] = await Promise.all([api.getEmployees()]);
  employeeSelect.innerHTML = `
    <option value="">Все сотрудники</option>
    ${(employeesResponse.employees || [])
      .map((employee) => `<option value="${employee.id}">${escapeHtml(employee.full_name)}</option>`)
      .join("")}
  `;

  filterForm.addEventListener("change", async () => {
    currentPage = 1;
    await loadOrders();
  });
  document.querySelector("#orders-filters-reset").addEventListener("click", async () => {
    filterForm.reset();
    currentPage = 1;
    await loadOrders();
  });
  document.querySelector("#create-order-button").addEventListener("click", () => openCreateOrderModal(loadOrders));
  await loadOrders();
};

const openOrderDetails = async (orderId) => {
  openModal(`<div class="status">Загрузка деталей заказа...</div>`);
  try {
    const response = await api.getOrder(orderId);
    const order = response.order;
    openModal(`
      <div class="section-head">
        <div>
          <h3>Заказ ${escapeHtml(order.id)}</h3>
          <p class="muted">${formatDateTime(order.created_at)}</p>
        </div>
        <button class="btn" id="close-modal-button">Закрыть</button>
      </div>
      <div class="grid two">
        <div class="card">
          <strong>Сотрудник</strong>
          <div>${escapeHtml(order.employee?.full_name || "—")}</div>
        </div>
        <div class="card">
          <strong>Итоговая сумма</strong>
          <div>${formatMoney(order.total_price)}</div>
        </div>
      </div>
      <div class="card">
        <h4>Позиции заказа</h4>
        <table>
          <thead>
            <tr>
              <th>Позиция</th>
              <th>Количество</th>
              <th>Цена</th>
              <th>Сумма</th>
            </tr>
          </thead>
          <tbody>
            ${(order.items || [])
              .map(
                (item) => `
                  <tr>
                    <td>${escapeHtml(item.menu_item_name)}</td>
                    <td>${item.quantity}</td>
                    <td>${formatMoney(item.price)}</td>
                    <td>${formatMoney(item.total_price)}</td>
                  </tr>
                `,
              )
              .join("")}
          </tbody>
        </table>
      </div>
      <div class="card">
        <h4>Списания ингредиентов</h4>
        ${
          order.inventory_operations?.length
            ? `
              <table>
                <thead>
                  <tr>
                    <th>Ингредиент</th>
                    <th>Изменение</th>
                    <th>Единица</th>
                    <th>Тип</th>
                  </tr>
                </thead>
                <tbody>
                  ${order.inventory_operations
                    .map(
                      (operation) => `
                        <tr>
                          <td>${escapeHtml(operation.ingredient_name)}</td>
                          <td>${formatNumber(operation.change_amount)}</td>
                          <td>${escapeHtml(operation.unit)}</td>
                          <td>${escapeHtml(formatOperationType(operation.operation_type))}</td>
                        </tr>
                      `,
                    )
                    .join("")}
                </tbody>
              </table>
            `
            : `<p class="muted">Для этого заказа складские списания не показаны.</p>`
        }
      </div>
    `);
    document.querySelector("#close-modal-button").addEventListener("click", closeModal);
  } catch (error) {
    closeModal();
    handleApiError(error);
  }
};

const openCreateOrderModal = async (onSuccess) => {
  openModal(`<div class="status">Загрузка меню...</div>`);
  try {
    const menuResponse = await api.getMenuItems(true);
    const menuItems = menuResponse.menu_items || [];
    let rowId = 0;

    const renderModal = () => {
      openModal(`
        <div class="section-head">
          <div>
            <h3>Создать заказ</h3>
            <p class="muted">Добавьте позиции заказа и укажите количество.</p>
          </div>
          <button class="btn" id="close-modal-button">Закрыть</button>
        </div>
        <form id="create-order-form" class="grid">
          <div id="order-items-list" class="list"></div>
          <div class="btn-row">
            <button type="button" class="btn" id="add-order-item-button">Добавить позицию</button>
            <button type="submit" class="btn primary">Сохранить заказ</button>
          </div>
        </form>
      `);
      document.querySelector("#close-modal-button").addEventListener("click", closeModal);
      document.querySelector("#add-order-item-button").addEventListener("click", addRow);
      document.querySelector("#create-order-form").addEventListener("submit", submitForm);
      if (!document.querySelector("#order-items-list").children.length) {
        addRow();
      }
    };

    const addRow = () => {
      rowId += 1;
      const row = document.createElement("div");
      row.className = "list-item";
      row.dataset.rowId = String(rowId);
      row.innerHTML = `
        <div class="form-grid">
          <div class="field">
            <label>Пункт меню</label>
            <select name="menu_item_id" required>
              <option value="">Выберите пункт меню</option>
              ${menuItems
                .map(
                  (item) => `
                    <option value="${item.id}" ${item.is_active ? "" : "disabled"}>
                      ${escapeHtml(item.name)}${item.is_active ? "" : " (неактивен)"}
                    </option>
                  `,
                )
                .join("")}
            </select>
          </div>
          <div class="field">
            <label>Количество</label>
            <input name="quantity" type="number" min="1" step="1" value="1" required />
          </div>
          <div class="field">
            <label>&nbsp;</label>
            <button type="button" class="btn" data-remove-row="${rowId}">Удалить позицию</button>
          </div>
        </div>
      `;
      document.querySelector("#order-items-list").appendChild(row);
      row.querySelector("[data-remove-row]").addEventListener("click", () => {
        row.remove();
      });
    };

    const submitForm = async (event) => {
      event.preventDefault();
      const rows = [...document.querySelectorAll("#order-items-list .list-item")];
      const items = rows
        .map((row) => ({
          menu_item_id: row.querySelector("[name='menu_item_id']").value,
          quantity: Number(row.querySelector("[name='quantity']").value),
        }))
        .filter((item) => item.menu_item_id);

      if (!items.length) {
        showToast("Добавьте хотя бы одну позицию заказа", "error");
        return;
      }

      const submitButton = event.currentTarget.querySelector("button[type='submit']");
      submitButton.disabled = true;
      try {
        await api.createOrder({ items });
        closeModal();
        showToast("Заказ создан", "success");
        await onSuccess();
      } catch (error) {
        handleApiError(error);
      } finally {
        submitButton.disabled = false;
      }
    };

    renderModal();
  } catch (error) {
    closeModal();
    handleApiError(error);
  }
};

const renderMenuPage = async () => {
  renderPageContent(`
    <div class="section-head">
      <div>
        <h3>Пункты меню</h3>
        <p class="muted">После создания доступны только изменение цены и активности.</p>
      </div>
      <button class="btn primary" id="create-menu-button">Создать пункт меню</button>
    </div>
    <div id="menu-list" class="card-grid"><div class="status">Загрузка меню...</div></div>
  `);

  const loadMenu = async () => {
    renderStatusBlock("#menu-list", "Загрузка меню...");
    try {
      const response = await api.getMenuItems(true);
      const items = response.menu_items || [];
      document.querySelector("#menu-list").innerHTML = items.length
        ? items
            .map(
              (item) => `
                <article class="card">
                  <div class="section-head">
                    <div>
                      <h4>${escapeHtml(item.name)}</h4>
                      <div>${formatMoney(item.price)}</div>
                    </div>
                    ${badge(item.is_active)}
                  </div>
                  <div>
                    <strong>Рецепт</strong>
                    ${
                      item.recipe?.length
                        ? `
                          <ul class="recipe-list">
                            ${item.recipe
                              .map(
                                (recipeItem) => `
                                  <li>
                                    ${escapeHtml(recipeItem.ingredient_name)}: ${formatNumber(recipeItem.amount)} ${escapeHtml(recipeItem.unit)}
                                    ${recipeItem.ingredient_is_active ? "" : "(неактивен)"}
                                  </li>
                                `,
                              )
                              .join("")}
                          </ul>
                        `
                        : `<p class="muted">Рецепт пуст.</p>`
                    }
                  </div>
                  <div class="btn-row">
                    <button class="btn" data-price-edit="${item.id}">Изменить цену</button>
                    <button class="btn" data-menu-toggle="${item.id}" data-next-active="${!item.is_active}">
                      ${item.is_active ? "Сделать неактивным" : "Сделать активным"}
                    </button>
                  </div>
                </article>
              `,
            )
            .join("")
        : createEmptyState("Меню пока пустое", "Добавьте первый пункт меню.");

      document.querySelectorAll("[data-price-edit]").forEach((button) => {
        button.addEventListener("click", () => openUpdatePriceModal(button.dataset.priceEdit, loadMenu));
      });
      document.querySelectorAll("[data-menu-toggle]").forEach((button) => {
        button.addEventListener("click", async () => {
          try {
            await api.setMenuActivity(button.dataset.menuToggle, button.dataset.nextActive === "true");
            await loadMenu();
          } catch (error) {
            handleApiError(error);
          }
        });
      });
    } catch (error) {
      handleApiError(error);
      renderStatusBlock("#menu-list", "Не удалось загрузить меню");
    }
  };

  document.querySelector("#create-menu-button").addEventListener("click", () => openCreateMenuModal(loadMenu));
  await loadMenu();
};

const openUpdatePriceModal = async (menuItemId, onSuccess) => {
  openModal(`<div class="status">Загрузка пункта меню...</div>`, { small: true });
  try {
    const response = await api.getMenuItem(menuItemId);
    const menuItem = response.menu_item;
    openModal(`
      <div class="section-head">
        <div>
          <h3>Изменить цену</h3>
          <p class="muted">${escapeHtml(menuItem.name)}</p>
        </div>
        <button class="btn" id="close-modal-button">Закрыть</button>
      </div>
      <form id="update-price-form" class="grid">
        <div class="field">
          <label for="new-price">Новая цена</label>
          <input id="new-price" type="number" min="0" step="0.01" value="${menuItem.price}" required />
        </div>
        <button class="btn primary" type="submit">Сохранить цену</button>
      </form>
    `, { small: true });
    document.querySelector("#close-modal-button").addEventListener("click", closeModal);
    document.querySelector("#update-price-form").addEventListener("submit", async (event) => {
      event.preventDefault();
      const button = event.currentTarget.querySelector("button[type='submit']");
      button.disabled = true;
      try {
        await api.updateMenuPrice(menuItemId, Number(document.querySelector("#new-price").value));
        closeModal();
        await onSuccess();
        showToast("Цена обновлена", "success");
      } catch (error) {
        handleApiError(error);
        button.disabled = false;
      }
    });
  } catch (error) {
    closeModal();
    handleApiError(error);
  }
};

const openCreateMenuModal = async (onSuccess) => {
  openModal(`<div class="status">Загрузка ингредиентов...</div>`);
  try {
    const response = await api.getIngredients(true);
    const ingredients = response.ingredients || [];
    let rowId = 0;

    openModal(`
      <div class="section-head">
        <div>
          <h3>Создать пункт меню</h3>
          <p class="muted">Неактивные ингредиенты доступны только для просмотра и заблокированы в выборе.</p>
        </div>
        <button class="btn" id="close-modal-button">Закрыть</button>
      </div>
      <form id="create-menu-form" class="grid">
        <div class="form-grid">
          <div class="field">
            <label for="menu-name">Название</label>
            <input id="menu-name" name="name" required />
          </div>
          <div class="field">
            <label for="menu-price">Цена</label>
            <input id="menu-price" name="price" type="number" min="0" step="0.01" required />
          </div>
        </div>
        <div>
          <div class="section-head">
            <h4>Рецепт</h4>
            <button type="button" class="btn" id="add-recipe-row-button">Добавить ингредиент</button>
          </div>
          <div id="recipe-list" class="list"></div>
        </div>
        <button class="btn primary" type="submit">Создать пункт меню</button>
      </form>
    `);

    document.querySelector("#close-modal-button").addEventListener("click", closeModal);
    document.querySelector("#add-recipe-row-button").addEventListener("click", addRecipeRow);
    document.querySelector("#create-menu-form").addEventListener("submit", submitForm);
    addRecipeRow();

    function addRecipeRow() {
      rowId += 1;
      const row = document.createElement("div");
      row.className = "list-item";
      row.innerHTML = `
        <div class="form-grid">
          <div class="field">
            <label>Ингредиент</label>
            <select name="ingredient_id" required>
              <option value="">Выберите ингредиент</option>
              ${ingredients
                .map(
                  (ingredient) => `
                    <option value="${ingredient.id}" ${ingredient.is_active ? "" : "disabled"}>
                      ${escapeHtml(ingredient.name)}${ingredient.is_active ? "" : " (неактивен)"}
                    </option>
                  `,
                )
                .join("")}
            </select>
          </div>
          <div class="field">
            <label>Количество</label>
            <input name="amount" type="number" min="0.01" step="0.01" required />
          </div>
          <div class="field">
            <label>&nbsp;</label>
            <button type="button" class="btn" data-remove>Удалить</button>
          </div>
        </div>
      `;
      row.querySelector("[data-remove]").addEventListener("click", () => row.remove());
      document.querySelector("#recipe-list").appendChild(row);
    }

    async function submitForm(event) {
      event.preventDefault();
      const form = event.currentTarget;
      const submitButton = form.querySelector("button[type='submit']");
      const recipe = [...document.querySelectorAll("#recipe-list .list-item")]
        .map((row) => ({
          ingredient_id: row.querySelector("[name='ingredient_id']").value,
          amount: Number(row.querySelector("[name='amount']").value),
        }))
        .filter((item) => item.ingredient_id);

      if (!recipe.length) {
        showToast("Добавьте хотя бы один ингредиент в рецепт", "error");
        return;
      }

      submitButton.disabled = true;
      try {
        await api.createMenuItem({
          name: form.elements.name.value,
          price: Number(form.elements.price.value),
          recipe,
        });
        closeModal();
        await onSuccess();
        showToast("Пункт меню создан", "success");
      } catch (error) {
        handleApiError(error);
      } finally {
        submitButton.disabled = false;
      }
    }
  } catch (error) {
    closeModal();
    handleApiError(error);
  }
};

const renderIngredientsPage = async () => {
  renderPageContent(`
    <div class="section-head">
      <div>
        <h3>Ингредиенты</h3>
        <p class="muted">Создавайте ингредиенты и контролируйте остатки на складе.</p>
      </div>
      <button class="btn primary" id="create-ingredient-button">Создать ингредиент</button>
    </div>
    <div id="ingredients-table-block" class="table-wrap"><div class="status">Загрузка ингредиентов...</div></div>
  `);

  const loadIngredients = async () => {
    renderStatusBlock("#ingredients-table-block", "Загрузка ингредиентов...");
    try {
      const response = await api.getIngredients(true);
      const ingredients = response.ingredients || [];
      document.querySelector("#ingredients-table-block").innerHTML = ingredients.length
        ? `
            <table>
              <thead>
                <tr>
                  <th>Название</th>
                  <th>Единица</th>
                  <th>Текущий остаток</th>
                  <th>Статус</th>
                  <th>Рецептов</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                ${ingredients
                  .map(
                    (ingredient) => `
                      <tr>
                        <td>${escapeHtml(ingredient.name)}</td>
                        <td>${escapeHtml(ingredient.unit)}</td>
                        <td>${formatNumber(ingredient.current_stock)}</td>
                        <td>${badge(ingredient.is_active)}</td>
                        <td>${ingredient.used_in_recipes_count || 0}</td>
                        <td>
                          <div class="btn-row">
                            <button class="btn" data-ingredient-details="${ingredient.id}">Подробнее</button>
                            <button class="btn" data-ingredient-toggle="${ingredient.id}" data-next-active="${!ingredient.is_active}">
                              ${ingredient.is_active ? "Сделать неактивным" : "Сделать активным"}
                            </button>
                          </div>
                        </td>
                      </tr>
                    `,
                  )
                  .join("")}
              </tbody>
            </table>
          `
        : createEmptyState("Ингредиентов пока нет", "Добавьте первый ингредиент для работы со складом и меню.");

      document.querySelectorAll("[data-ingredient-details]").forEach((button) => {
        button.addEventListener("click", () => openIngredientDetails(button.dataset.ingredientDetails));
      });
      document.querySelectorAll("[data-ingredient-toggle]").forEach((button) => {
        button.addEventListener("click", async () => {
          try {
            await api.setIngredientActivity(
              button.dataset.ingredientToggle,
              button.dataset.nextActive === "true",
            );
            await loadIngredients();
          } catch (error) {
            if (error instanceof ApiError && error.code === "INGREDIENT_USED_IN_RECIPE") {
              const details = Array.isArray(error.details) ? error.details.join(", ") : "";
              showToast(
                details
                  ? `Ингредиент нельзя деактивировать. Используется в: ${details}`
                  : "Ингредиент нельзя сделать неактивным, потому что он используется в рецептах",
                "error",
              );
            } else {
              handleApiError(error);
            }
          }
        });
      });
    } catch (error) {
      handleApiError(error);
      renderStatusBlock("#ingredients-table-block", "Не удалось загрузить ингредиенты");
    }
  };

  document
    .querySelector("#create-ingredient-button")
    .addEventListener("click", () => openCreateIngredientModal(loadIngredients));
  await loadIngredients();
};

const openCreateIngredientModal = (onSuccess) => {
  openModal(`
    <div class="section-head">
      <div>
        <h3>Создать ингредиент</h3>
        <p class="muted">Поле initial_stock можно оставить пустым.</p>
      </div>
      <button class="btn" id="close-modal-button">Закрыть</button>
    </div>
    <form id="create-ingredient-form" class="grid">
      <div class="form-grid">
        <div class="field">
          <label for="ingredient-name">Название</label>
          <input id="ingredient-name" name="name" required />
        </div>
        <div class="field">
          <label for="ingredient-unit">Единица измерения</label>
          <input id="ingredient-unit" name="unit" required />
        </div>
        <div class="field">
          <label for="ingredient-stock">Начальный остаток</label>
          <input id="ingredient-stock" name="initial_stock" type="number" min="0" step="0.01" />
        </div>
      </div>
      <button class="btn primary" type="submit">Создать ингредиент</button>
    </form>
  `, { small: true });
  document.querySelector("#close-modal-button").addEventListener("click", closeModal);
  document.querySelector("#create-ingredient-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const form = event.currentTarget;
    const button = form.querySelector("button[type='submit']");
    button.disabled = true;
    try {
      await api.createIngredient({
        name: form.elements.name.value,
        unit: form.elements.unit.value,
        initial_stock: form.elements.initial_stock.value ? Number(form.elements.initial_stock.value) : null,
      });
      closeModal();
      await onSuccess();
      showToast("Ингредиент создан", "success");
    } catch (error) {
      handleApiError(error);
    } finally {
      button.disabled = false;
    }
  });
};

const openIngredientDetails = async (ingredientId) => {
  openModal(`<div class="status">Загрузка ингредиента...</div>`);
  try {
    const response = await api.getIngredient(ingredientId);
    const ingredient = response.ingredient;
    openModal(`
      <div class="section-head">
        <div>
          <h3>${escapeHtml(ingredient.name)}</h3>
          <p class="muted">Подробности ингредиента</p>
        </div>
        <button class="btn" id="close-modal-button">Закрыть</button>
      </div>
      <div class="grid two">
        <div class="card">
          <strong>Единица</strong>
          <div>${escapeHtml(ingredient.unit)}</div>
        </div>
        <div class="card">
          <strong>Остаток</strong>
          <div>${formatNumber(ingredient.current_stock)}</div>
        </div>
      </div>
      <div class="card">
        <strong>Статус</strong>
        <div>${badge(ingredient.is_active)}</div>
      </div>
      <div class="card">
        <h4>Использование в рецептах</h4>
        ${
          ingredient.used_in_recipes?.length
            ? `
              <ul class="recipe-list">
                ${ingredient.used_in_recipes
                  .map(
                    (usage) => `
                      <li>${escapeHtml(usage.menu_item_name)} — ${formatNumber(usage.amount)}</li>
                    `,
                  )
                  .join("")}
              </ul>
            `
            : `<p class="muted">Ингредиент не используется в рецептах.</p>`
        }
      </div>
    `);
    document.querySelector("#close-modal-button").addEventListener("click", closeModal);
  } catch (error) {
    closeModal();
    handleApiError(error);
  }
};

const renderInventoryPage = async () => {
  let currentPage = 1;

  renderPageContent(`
    <div class="section-head">
      <div>
        <h3>Складские операции</h3>
        <p class="muted">Просматривайте движения по складу и добавляйте ручные операции.</p>
      </div>
      <button class="btn primary" id="create-operation-button">Создать операцию</button>
    </div>
    <form id="inventory-filter-form" class="filter-bar">
      <div class="form-grid">
        <div class="field">
          <label>Сотрудник</label>
          <select name="employee_id" id="inventory-employee-filter"></select>
        </div>
        <div class="field">
          <label>Ингредиент</label>
          <select name="ingredient_id" id="inventory-ingredient-filter"></select>
        </div>
        <div class="field">
          <label>Тип операции</label>
          <select name="operation_type" id="inventory-type-filter"></select>
        </div>
        <div class="field">
          <label>from</label>
          <input type="datetime-local" name="from" />
        </div>
        <div class="field">
          <label>to</label>
          <input type="datetime-local" name="to" />
        </div>
      </div>
      <div class="btn-row">
        <button type="button" class="btn" id="inventory-filters-reset">Очистить фильтры</button>
      </div>
    </form>
    <div id="inventory-table-block" class="table-wrap"><div class="status">Загрузка операций...</div></div>
  `);

  const filterForm = document.querySelector("#inventory-filter-form");

  const [employeesResponse, ingredientsResponse] = await Promise.all([
    api.getEmployees(),
    api.getIngredients(true),
  ]);

  document.querySelector("#inventory-employee-filter").innerHTML = `
    <option value="">Все сотрудники</option>
    ${(employeesResponse.employees || [])
      .map((employee) => `<option value="${employee.id}">${escapeHtml(employee.full_name)}</option>`)
      .join("")}
  `;

  document.querySelector("#inventory-ingredient-filter").innerHTML = `
    <option value="">Все ингредиенты</option>
    ${(ingredientsResponse.ingredients || [])
      .map((ingredient) => `<option value="${ingredient.id}">${escapeHtml(ingredient.name)}</option>`)
      .join("")}
  `;

  document.querySelector("#inventory-type-filter").innerHTML = `
    <option value="">Все типы</option>
    ${OPERATION_TYPE_OPTIONS.map((option) => `<option value="${option.value}">${option.label}</option>`).join("")}
  `;

  const loadOperations = async () => {
    renderStatusBlock("#inventory-table-block", "Загрузка операций...");
    try {
      const response = await api.getOperations({
        page: currentPage,
        employee_id: filterForm.elements.employee_id.value,
        ingredient_id: filterForm.elements.ingredient_id.value,
        operation_type: filterForm.elements.operation_type.value,
        from: filterForm.elements.from.value ? new Date(filterForm.elements.from.value).toISOString() : "",
        to: filterForm.elements.to.value ? new Date(filterForm.elements.to.value).toISOString() : "",
      });
      const operations = response.items || [];
      document.querySelector("#inventory-table-block").innerHTML = operations.length
        ? `
            <table>
              <thead>
                <tr>
                  <th>Дата</th>
                  <th>Тип операции</th>
                  <th>Ингредиент</th>
                  <th>Изменение</th>
                  <th>Единица</th>
                  <th>Сотрудник</th>
                  <th>Заказ</th>
                </tr>
              </thead>
              <tbody>
                ${operations
                  .map(
                    (operation) => `
                      <tr>
                        <td>${formatDateTime(operation.created_at)}</td>
                        <td>${escapeHtml(formatOperationType(operation.operation_type))}</td>
                        <td>${escapeHtml(operation.ingredient?.name || "—")}</td>
                        <td>${formatNumber(operation.change_amount)}</td>
                        <td>${escapeHtml(operation.ingredient?.unit || "—")}</td>
                        <td>${escapeHtml(operation.employee?.full_name || "—")}</td>
                        <td>${escapeHtml(operation.order?.id || "—")}</td>
                      </tr>
                    `,
                  )
                  .join("")}
              </tbody>
            </table>
            ${buildPaginationInfo(response.pagination)}
            ${buildPaginationControls(response.pagination, "inventory")}
          `
        : createEmptyState("Операций пока нет", "Создайте ручную операцию или дождитесь автоматических списаний по заказам.");

      document.querySelectorAll('[data-pagination-prefix="inventory"]').forEach((button) => {
        button.addEventListener("click", async () => {
          currentPage = Number(button.dataset.page);
          await loadOperations();
        });
      });
    } catch (error) {
      handleApiError(error);
      renderStatusBlock("#inventory-table-block", "Не удалось загрузить операции");
    }
  };

  filterForm.addEventListener("change", async () => {
    currentPage = 1;
    await loadOperations();
  });
  document.querySelector("#inventory-filters-reset").addEventListener("click", async () => {
    filterForm.reset();
    currentPage = 1;
    await loadOperations();
  });
  document.querySelector("#create-operation-button").addEventListener("click", () => {
    openCreateOperationModal(ingredientsResponse.ingredients || [], async () => {
      await loadOperations();
    });
  });
  await loadOperations();
};

const openCreateOperationModal = (ingredients, onSuccess) => {
  openModal(`
    <div class="section-head">
      <div>
        <h3>Создать складскую операцию</h3>
        <p class="muted">Для создания доступны только ручные типы и активные ингредиенты.</p>
      </div>
      <button class="btn" id="close-modal-button">Закрыть</button>
    </div>
    <form id="create-operation-form" class="grid">
      <div class="form-grid">
        <div class="field">
          <label>Ингредиент</label>
          <select name="ingredient_id" required>
            <option value="">Выберите ингредиент</option>
            ${ingredients
              .map(
                (ingredient) => `
                  <option value="${ingredient.id}" ${ingredient.is_active ? "" : "disabled"}>
                    ${escapeHtml(ingredient.name)}${ingredient.is_active ? "" : " (неактивен)"}
                  </option>
                `,
              )
              .join("")}
          </select>
        </div>
        <div class="field">
          <label>Тип операции</label>
          <select name="operation_type" required>
            <option value="">Выберите тип</option>
            ${MANUAL_OPERATION_TYPES.map(
              (option) => `<option value="${option.value}">${option.label}</option>`,
            ).join("")}
          </select>
        </div>
        <div class="field">
          <label>Количество</label>
          <input name="amount" type="number" min="0.01" step="0.01" required />
        </div>
      </div>
      <button class="btn primary" type="submit">Создать операцию</button>
    </form>
  `, { small: true });

  document.querySelector("#close-modal-button").addEventListener("click", closeModal);
  document.querySelector("#create-operation-form").addEventListener("submit", async (event) => {
    event.preventDefault();
    const form = event.currentTarget;
    const button = form.querySelector("button[type='submit']");
    button.disabled = true;
    try {
      await api.createOperation({
        ingredient_id: form.elements.ingredient_id.value,
        operation_type: form.elements.operation_type.value,
        amount: Number(form.elements.amount.value),
      });
      closeModal();
      await onSuccess();
      showToast("Операция создана", "success");
    } catch (error) {
      if (error instanceof ApiError && error.code === "INSUFFICIENT_STOCK") {
        showToast("Недостаточно остатка для списания", "error");
      } else {
        handleApiError(error);
      }
    } finally {
      button.disabled = false;
    }
  });
};

const renderReportsPage = async () => {
  renderPageContent(`
    <div class="section-head">
      <div>
        <h3>PDF-отчеты</h3>
        <p class="muted">Выберите период и скачайте нужный отчет в PDF.</p>
      </div>
    </div>
    <div id="reports-grid" class="card-grid"><div class="status">Загрузка справочников...</div></div>
  `);

  try {
    const [employeesResponse, ingredientsResponse] = await Promise.all([
      api.getEmployees(),
      api.getIngredients(true),
    ]);

    const employeeOptions = `
      <option value="">Все сотрудники</option>
      ${(employeesResponse.employees || [])
        .map((employee) => `<option value="${employee.id}">${escapeHtml(employee.full_name)}</option>`)
        .join("")}
    `;
    const ingredientOptions = `
      <option value="">Все ингредиенты</option>
      ${(ingredientsResponse.ingredients || [])
        .map((ingredient) => `<option value="${ingredient.id}">${escapeHtml(ingredient.name)}</option>`)
        .join("")}
    `;
    const defaults = getDefaultDateRange();

    document.querySelector("#reports-grid").innerHTML = `
      ${reportCard("sales", "Отчет по продажам", employeeOptions, ingredientOptions, defaults)}
      ${reportCard("employees", "Отчет по сотрудникам", employeeOptions, ingredientOptions, defaults)}
      ${reportCard("inventory", "Отчет по складу", employeeOptions, ingredientOptions, defaults, true)}
      ${reportCard("popularity", "Отчет по популярности меню", employeeOptions, ingredientOptions, defaults)}
    `;

    document.querySelectorAll("[data-report-form]").forEach((form) => {
      form.addEventListener("submit", async (event) => {
        event.preventDefault();
        await submitReportForm(event.currentTarget);
      });
    });
  } catch (error) {
    handleApiError(error);
    renderStatusBlock("#reports-grid", "Не удалось загрузить данные для отчетов");
  }
};

const reportCard = (type, title, employeeOptions, ingredientOptions, defaults, isInventory = false) => `
  <section class="card">
    <h3>${title}</h3>
    <form class="grid" data-report-form="${type}">
      <input type="hidden" name="report_type" value="${type}" />
      <div class="field">
        <label>from</label>
        <input type="datetime-local" name="from" value="${defaults.from}" required />
      </div>
      <div class="field">
        <label>to</label>
        <input type="datetime-local" name="to" value="${defaults.to}" required />
      </div>
      <div class="field">
        <label>Сотрудник</label>
        <select name="employee_id">${employeeOptions}</select>
      </div>
      ${
        isInventory
          ? `
            <div class="field">
              <label>Ингредиент</label>
              <select name="ingredient_id">${ingredientOptions}</select>
            </div>
            <div class="field">
              <label>Тип операции</label>
              <select name="operation_type">
                <option value="">Все типы</option>
                ${OPERATION_TYPE_OPTIONS.map(
                  (option) => `<option value="${option.value}">${option.label}</option>`,
                ).join("")}
              </select>
            </div>
          `
          : ""
      }
      <button class="btn primary" type="submit">Скачать PDF</button>
    </form>
  </section>
`;

const submitReportForm = async (form) => {
  const reportType = form.elements.report_type.value;
  const payload = {
    from: new Date(form.elements.from.value).toISOString(),
    to: new Date(form.elements.to.value).toISOString(),
    employee_id: form.elements.employee_id.value || null,
  };

  if (reportType === "inventory") {
    payload.ingredient_id = form.elements.ingredient_id.value || null;
    payload.operation_type = form.elements.operation_type.value || null;
  }

  const button = form.querySelector("button[type='submit']");
  button.disabled = true;
  try {
    const blob =
      reportType === "sales"
        ? await api.downloadSalesReport(payload)
        : reportType === "employees"
          ? await api.downloadEmployeesReport(payload)
          : reportType === "inventory"
            ? await api.downloadInventoryReport(payload)
            : await api.downloadPopularityReport(payload);

    downloadBlob(blob, REPORT_FILES[reportType]);
    showToast("Отчет скачан", "success");
  } catch (error) {
    handleApiError(error);
  } finally {
    button.disabled = false;
  }
};

const renderCurrentProtectedPage = async () => {
  renderProtectedLayout();
  const path = window.location.pathname;
  if (path === ROUTES.orders) {
    await renderOrdersPage();
  } else if (path === ROUTES.menu) {
    await renderMenuPage();
  } else if (path === ROUTES.ingredients) {
    await renderIngredientsPage();
  } else if (path === ROUTES.inventory) {
    await renderInventoryPage();
  } else if (path === ROUTES.reports) {
    await renderReportsPage();
  } else {
    navigate(ROUTES.orders, { replace: true });
  }
};

const renderApp = async () => {
  const path = window.location.pathname;
  renderLoading();
  try {
    if (path === ROUTES.shift) {
      const active = await api.getActiveShift();
      if (active.is_active) {
        state.activeShift = active.shift;
        navigate(ROUTES.orders, { replace: true });
        return;
      }
      state.activeShift = null;
      renderOpenShiftScreen();
      return;
    }

    const isShiftActive = await ensureActiveShift();
    if (isShiftActive) {
      await renderCurrentProtectedPage();
    }
  } catch (error) {
    if (path !== ROUTES.shift) {
      handleApiError(error);
    } else {
      showToast(getErrorMessage(error), "error");
      renderOpenShiftScreen();
    }
  }
};

window.addEventListener("popstate", renderApp);
window.addEventListener("DOMContentLoaded", renderApp);
