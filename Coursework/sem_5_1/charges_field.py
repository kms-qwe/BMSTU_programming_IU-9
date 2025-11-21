import math
import pyglet

# ----- ВАЖНО ДЛЯ RETINA / HiDPI -----
pyglet.options.dpi_scaling = "stretch"

from pyglet import shapes

# ================== НАСТРАИВАЕМЫЕ ПАРАМЕТРЫ ==================

# Границы видимой области по x и y: [-L, L]
L = 10.0

# Размер сетки для расчёта потенциала (N x N).
N_GRID = 220

# Размер квадратной области поля в логических пикселях
WINDOW_FIELD_SIZE = 800

# Ширина панели с легендой справа, в пикселях
LEGEND_WIDTH = 120

# Итоговый размер окна
WINDOW_WIDTH = WINDOW_FIELD_SIZE + LEGEND_WIDTH
WINDOW_HEIGHT = WINDOW_FIELD_SIZE

# Физические параметры
R_MIN = 0.05             # минимальный радиус для поля/силы
MERGE_THRESHOLD = 0.25   # порог расстояния для слияния зарядов
PHYSICS_DT = 0.01        # шаг интегрирования по времени
DAMPING = 0.98           # демпфирование скоростей
K_COULOMB = 1.0          # константа в законе Кулона

# Период пересчёта поля/потенциала и линий напряжённости (сек)
FIELD_UPDATE_INTERVAL = 0.08

# Параметры интегрирования линий поля
STREAMLINE_STEP = 0.04        # шаг по полю (мировые координаты)
STREAMLINE_MAX_STEPS = 600    # максимум шагов на линию

# Радиус, на котором линия напряжённости считается дошедшей до отрицательного заряда
FIELDLINE_END_RADIUS = 0.08

# Ограничение динамического диапазона для раскраски потенциала.
POTENTIAL_COLOR_MAX = 4.0

# Список исходных зарядов: (x, y, q)
INITIAL_CHARGES = [
    (-1.5, 0.0, +2.0),
    (1.5, 0.0, -2.0),
    (0.0, 1.5, +1.0),
    (0.0, -1.5, -1.0),
    (-3.0, -3.5, -1.0),
    (4.0, 3.5, +2.0)

]


class Charge:
    """Класс точечного заряда в мировых координатах."""

    def __init__(self, x, y, q, vx=0.0, vy=0.0, m=1.0):
        self.x = float(x)
        self.y = float(y)
        self.vx = float(vx)
        self.vy = float(vy)
        self.q = float(q)
        self.m = float(m)
        self.shape: shapes.Circle | None = None

    def apply_force(self, fx: float, fy: float, dt: float, damping: float) -> None:
        ax = fx / self.m
        ay = fy / self.m
        self.vx += ax * dt
        self.vy += ay * dt
        self.vx *= damping
        self.vy *= damping

    def update_position(self, dt: float) -> None:
        self.x += self.vx * dt
        self.y += self.vy * dt


class FieldSimulationWindow(pyglet.window.Window):
    """Основное окно симуляции: физика + отрисовка."""

    def __init__(self):
        super().__init__(
            width=WINDOW_WIDTH,
            height=WINDOW_HEIGHT,
            caption="Визуализация электрических полей",
            resizable=False,
        )

        # тёмный фон вокруг поля
        pyglet.gl.glClearColor(0.06, 0.06, 0.08, 1.0)

        self.field_pixel_size = WINDOW_FIELD_SIZE
        self.legend_width = LEGEND_WIDTH
        self.L = L

        # Один общий Batch + слои (Group) для порядка отрисовки
        self.batch = pyglet.graphics.Batch()
        self.group_background = pyglet.graphics.Group(order=0)   # градиент потенциала
        self.group_fieldlines = pyglet.graphics.Group(order=1)   # линии напряжённости
        self.group_axes = pyglet.graphics.Group(order=2)         # оси, насечки, подписи
        self.group_charges = pyglet.graphics.Group(order=3)      # заряды
        self.group_legend = pyglet.graphics.Group(order=4)       # легенда

        # Ссылка на все объекты осей (чтоб их не съел GC)
        self.axes_drawables: list[object] = []

        # ---------- Заряды ----------
        self.charges: list[Charge] = []
        for x, y, q in INITIAL_CHARGES:
            self.add_charge(x, y, q)

        # ---------- Сетка потенциальных точек (для раскраски) ----------
        self.grid_resolution = N_GRID
        self.grid_dx = (2.0 * self.L) / self.grid_resolution
        self.x_world = [
            -self.L + (i + 0.5) * self.grid_dx for i in range(self.grid_resolution)
        ]
        self.y_world = [
            -self.L + (j + 0.5) * self.grid_dx for j in range(self.grid_resolution)
        ]

        self.potential_values: list[float] = [
            0.0 for _ in range(self.grid_resolution * self.grid_resolution)
        ]
        self.potential_cells: list[shapes.Rectangle] = []
        self.global_vmax = 1.0  # диапазон для легенды

        self._create_potential_cells()
        self._create_axes()
        self._create_legend_shapes()

        # Линии поля (streamlines)
        self.streamline_shapes: list[shapes.MultiLine] = []

        # Таймеры
        self.time_accumulator = 0.0
        self.time_since_field_update = 0.0

        pyglet.clock.schedule_interval(self.update_simulation, 1.0 / 60.0)

    # =========================================================
    #           КООРДИНАТЫ
    # =========================================================

    def world_to_screen(self, x: float, y: float) -> tuple[float, float]:
        sx = (x + self.L) / (2.0 * self.L) * self.field_pixel_size
        sy = (y + self.L) / (2.0 * self.L) * self.field_pixel_size
        return sx, sy

    def screen_to_world(self, sx: float, sy: float) -> tuple[float, float]:
        x = sx / self.field_pixel_size * (2.0 * self.L) - self.L
        y = sy / self.field_pixel_size * (2.0 * self.L) - self.L
        return x, y

    # =========================================================
    #           СОЗДАНИЕ ОБЪЕКТОВ
    # =========================================================

    def add_charge(self, x: float, y: float, q: float) -> None:
        charge = Charge(x, y, q)
        radius = self._charge_screen_radius(q)
        color = (255, 80, 80) if q > 0 else (80, 80, 255)
        sx, sy = self.world_to_screen(charge.x, charge.y)
        circle = shapes.Circle(
            sx,
            sy,
            radius,
            segments=None,
            color=color,
            batch=self.batch,
            group=self.group_charges,
        )
        charge.shape = circle
        self.charges.append(charge)

    def _charge_screen_radius(self, q: float) -> float:
        base = 8.0
        return base * (0.7 + 0.3 * abs(q))

    def _create_potential_cells(self) -> None:
        """Создаём прямоугольники для градиента потенциала."""
        cell_w = self.field_pixel_size / self.grid_resolution
        cell_h = self.field_pixel_size / self.grid_resolution
        for j in range(self.grid_resolution):
            for i in range(self.grid_resolution):
                x_px = i * cell_w
                y_px = j * cell_h
                rect = shapes.Rectangle(
                    x_px,
                    y_px,
                    cell_w,
                    cell_h,
                    color=(0, 0, 0),
                    batch=self.batch,
                    group=self.group_background,
                )
                rect.opacity = 235
                self.potential_cells.append(rect)

    def _create_axes(self) -> None:
        """Создаём оси OX и OY, насечки и подписи."""
        label_values = [-self.L, -self.L / 2.0, 0.0, self.L / 2.0, self.L]
        tick_size = 10.0

        # ----- ОСИ X и Y -----
        # Ось X: y = 0
        sx1, sy1 = self.world_to_screen(-self.L, 0.0)
        sx2, sy2 = self.world_to_screen(self.L, 0.0)
        axis_x = shapes.Line(
            sx1,
            sy1,
            sx2,
            sy2,
            thickness=3.0,
            color=(0, 0, 0),
            batch=self.batch,
            group=self.group_axes,
        )
        self.axes_drawables.append(axis_x)

        # Ось Y: x = 0
        sx1, sy1 = self.world_to_screen(0.0, -self.L)
        sx2, sy2 = self.world_to_screen(0.0, self.L)
        axis_y = shapes.Line(
            sx1,
            sy1,
            sx2,
            sy2,
            thickness=3.0,
            color=(0, 0, 0),
            batch=self.batch,
            group=self.group_axes,
        )
        self.axes_drawables.append(axis_y)

        # ----- НАСЕЧКИ И ПОДПИСИ -----

        # X-ось: вертикальные насечки
        for xv in label_values:
            sx, sy = self.world_to_screen(xv, 0.0)
            tick = shapes.Line(
                sx,
                sy - tick_size / 2.0,
                sx,
                sy + tick_size / 2.0,
                thickness=2.0,
                color=(0, 0, 0),
                batch=self.batch,
                group=self.group_axes,
            )
            self.axes_drawables.append(tick)

        # Числовые подписи для оси X — ПОД осью
        for xv in label_values:
            sx, sy = self.world_to_screen(xv, 0.0)
            label = pyglet.text.Label(
                f"{xv:.1f}",
                font_size=11,
                x=sx,
                y=sy - tick_size - 4,  # чуть ниже оси
                anchor_x="center",
                anchor_y="top",
                color=(0, 0, 0, 255),
                batch=self.batch,
                group=self.group_axes,
            )
            self.axes_drawables.append(label)

        # Y-ось: горизонтальные насечки
        for yv in label_values:
            sx, sy = self.world_to_screen(0.0, yv)
            tick = shapes.Line(
                sx - tick_size / 2.0,
                sy,
                sx + tick_size / 2.0,
                sy,
                thickness=2.0,
                color=(0, 0, 0),
                batch=self.batch,
                group=self.group_axes,
            )
            self.axes_drawables.append(tick)

        # Числовые подписи для оси Y — СПРАВА от оси
        for yv in label_values:
            sx, sy = self.world_to_screen(0.0, yv)
            label = pyglet.text.Label(
                f"{yv:.1f}",
                font_size=11,
                x=sx + tick_size + 4,  # чуть правее оси
                y=sy,
                anchor_x="left",
                anchor_y="center",
                color=(0, 0, 0, 255),
                batch=self.batch,
                group=self.group_axes,
            )
            self.axes_drawables.append(label)

        # Метка оси X
        sx_x, sy_x = self.world_to_screen(self.L, 0.0)
        label_x = pyglet.text.Label(
            "x",
            font_size=14,
            x=sx_x + 18,
            y=sy_x + 8,
            anchor_x="left",
            anchor_y="bottom",
            color=(0, 0, 0, 255),
            batch=self.batch,
            group=self.group_axes,
        )
        self.axes_drawables.append(label_x)

        # Метка оси Y
        sx_y, sy_y = self.world_to_screen(0.0, self.L)
        label_y = pyglet.text.Label(
            "y",
            font_size=14,
            x=sx_y + 6,
            y=sy_y + 18,
            anchor_x="left",
            anchor_y="bottom",
            color=(0, 0, 0, 255),
            batch=self.batch,
            group=self.group_axes,
        )
        self.axes_drawables.append(label_y)

    def _create_legend_shapes(self) -> None:
        bar_width = 40
        margin_x = 20
        margin_y = 60
        bar_x = self.field_pixel_size + margin_x
        bar_height = self.height - 2 * margin_y
        bar_y = margin_y
        num_segments = 60

        self.legend_rectangles: list[shapes.Rectangle] = []
        segment_h = bar_height / num_segments

        for k in range(num_segments):
            v_norm = -1.0 + 2.0 * (k + 0.5) / num_segments
            fake_V = v_norm
            r, g, b = self._potential_to_color(fake_V, 1.0)
            rect = shapes.Rectangle(
                bar_x,
                bar_y + k * segment_h,
                bar_width,
                segment_h + 1,
                color=(r, g, b),
                batch=self.batch,
                group=self.group_legend,
            )
            self.legend_rectangles.append(rect)

        shapes.Box(
            bar_x,
            bar_y,
            bar_width,
            bar_height,
            thickness=1.5,
            color=(230, 230, 230),
            batch=self.batch,
            group=self.group_legend,
        )

        y_top = bar_y + bar_height
        y_mid = bar_y + bar_height / 2.0
        y_mid_pos = bar_y + bar_height * 0.75
        y_mid_neg = bar_y + bar_height * 0.25

        tick_len = 6
        for y_tick in (y_top, y_mid_pos, y_mid, y_mid_neg, bar_y):
            shapes.Line(
                bar_x + bar_width,
                y_tick,
                bar_x + bar_width + tick_len,
                y_tick,
                thickness=1.0,
                color=(230, 230, 230),
                batch=self.batch,
                group=self.group_legend,
            )

        self.legend_label_top = pyglet.text.Label(
            "+Vmax",
            font_size=10,
            x=bar_x + bar_width + tick_len + 4,
            y=y_top,
            anchor_x="left",
            anchor_y="top",
            color=(230, 230, 230, 255),
            batch=self.batch,
            group=self.group_legend,
        )
        self.legend_label_mid_pos = pyglet.text.Label(
            "+Vmax/2",
            font_size=10,
            x=bar_x + bar_width + tick_len + 4,
            y=y_mid_pos,
            anchor_x="left",
            anchor_y="center",
            color=(230, 230, 230, 255),
            batch=self.batch,
            group=self.group_legend,
        )
        self.legend_label_mid = pyglet.text.Label(
            "0",
            font_size=10,
            x=bar_x + bar_width + tick_len + 4,
            y=y_mid,
            anchor_x="left",
            anchor_y="center",
            color=(230, 230, 230, 255),
            batch=self.batch,
            group=self.group_legend,
        )
        self.legend_label_mid_neg = pyglet.text.Label(
            "-Vmax/2",
            font_size=10,
            x=bar_x + bar_width + tick_len + 4,
            y=y_mid_neg,
            anchor_x="left",
            anchor_y="center",
            color=(230, 230, 230, 255),
            batch=self.batch,
            group=self.group_legend,
        )
        self.legend_label_bot = pyglet.text.Label(
            "-Vmax",
            font_size=10,
            x=bar_x + bar_width + tick_len + 4,
            y=bar_y,
            anchor_x="left",
            anchor_y="bottom",
            color=(230, 230, 230, 255),
            batch=self.batch,
            group=self.group_legend,
        )

        self.legend_title = pyglet.text.Label(
            "Potential V(x,y)",
            font_size=11,
            x=bar_x + bar_width / 2.0 + 10,
            y=bar_y + bar_height + 25,
            anchor_x="center",
            anchor_y="bottom",
            color=(240, 240, 240, 255),
            batch=self.batch,
            group=self.group_legend,
        )

    # =========================================================
    #           ФИЗИКА
    # =========================================================

    def compute_forces(self) -> list[tuple[float, float]]:
        n = len(self.charges)
        forces: list[tuple[float, float]] = [(0.0, 0.0) for _ in range(n)]
        for i in range(n):
            ci = self.charges[i]
            fx_i = 0.0
            fy_i = 0.0
            for j in range(n):
                if i == j:
                    continue
                cj = self.charges[j]
                dx = ci.x - cj.x
                dy = ci.y - cj.y
                r2 = dx * dx + dy * dy
                if r2 < R_MIN * R_MIN:
                    r2 = R_MIN * R_MIN
                r = math.sqrt(r2)
                if r == 0.0:
                    continue
                force_mag = K_COULOMB * ci.q * cj.q / r2
                fx_i += force_mag * dx / r
                fy_i += force_mag * dy / r
            forces[i] = (fx_i, fy_i)
        return forces

    def merge_charges(self) -> None:
        merged = True
        while merged:
            merged = False
            i = 0
            while i < len(self.charges):
                j = i + 1
                while j < len(self.charges):
                    ci = self.charges[i]
                    cj = self.charges[j]
                    if ci.q * cj.q < 0.0:
                        dx = ci.x - cj.x
                        dy = ci.y - cj.y
                        dist2 = dx * dx + dy * dy
                        if dist2 < MERGE_THRESHOLD * MERGE_THRESHOLD:
                            q_new = ci.q + cj.q
                            x_new = 0.5 * (ci.x + cj.x)
                            y_new = 0.5 * (ci.y + cj.y)
                            vx_new = 0.5 * (ci.vx + cj.vx)
                            vy_new = 0.5 * (ci.vy + cj.vy)

                            ci.shape.delete()
                            cj.shape.delete()
                            self.charges.pop(j)
                            self.charges.pop(i)

                            if abs(q_new) > 1e-6:
                                new_charge = Charge(x_new, y_new, q_new, vx_new, vy_new)
                                radius = self._charge_screen_radius(q_new)
                                color = (255, 80, 80) if q_new > 0 else (80, 80, 255)
                                sx, sy = self.world_to_screen(x_new, y_new)
                                circle = shapes.Circle(
                                    sx,
                                    sy,
                                    radius,
                                    segments=None,
                                    color=color,
                                    batch=self.batch,
                                    group=self.group_charges,
                                )
                                new_charge.shape = circle
                                self.charges.insert(i, new_charge)

                            merged = True
                            i = -1
                            break
                    j += 1
                i += 1

    def compute_potential_grid(self) -> None:
        local_max = 0.0
        idx = 0
        for yw in self.y_world:
            for xw in self.x_world:
                V = 0.0
                for c in self.charges:
                    dx = xw - c.x
                    dy = yw - c.y
                    r2 = dx * dx + dy * dy
                    if r2 < R_MIN * R_MIN:
                        r2 = R_MIN * R_MIN
                    r = math.sqrt(r2)
                    V += c.q / r
                self.potential_values[idx] = V
                if abs(V) > local_max:
                    local_max = abs(V)
                idx += 1

        if local_max < 1e-6:
            local_max = 1.0
        display_max = min(local_max, POTENTIAL_COLOR_MAX)
        self.global_vmax = display_max

        for idx, V in enumerate(self.potential_values):
            r, g, b = self._potential_to_color(V, self.global_vmax)
            self.potential_cells[idx].color = (r, g, b)

        vmax = self.global_vmax
        self.legend_label_top.text = f"+{vmax:.2f}"
        self.legend_label_mid_pos.text = f"+{0.5 * vmax:.2f}"
        self.legend_label_mid.text = "0"
        self.legend_label_mid_neg.text = f"-{0.5 * vmax:.2f}"
        self.legend_label_bot.text = f"-{vmax:.2f}"

        num_segments = len(self.legend_rectangles)
        for k, rect in enumerate(self.legend_rectangles):
            v_norm = -1.0 + 2.0 * (k + 0.5) / num_segments
            fake_V = v_norm * self.global_vmax
            r, g, b = self._potential_to_color(fake_V, self.global_vmax)
            rect.color = (r, g, b)

    def _potential_to_color(self, V: float, Vmax: float) -> tuple[int, int, int]:
        if Vmax <= 0.0:
            Vmax = 1.0
        v = max(-Vmax, min(Vmax, V)) / Vmax
        levels = 40
        v = round(v * levels) / levels

        if v >= 0.0:
            t = v
            r = 255
            g = int(255 * (1.0 - 0.6 * t))  # 255 -> 150
            b = int(255 * (1.0 - t))        # 255 -> 0
        else:
            t = -v
            r = int(255 * (1.0 - t))        # 255 -> 0
            g = int(255 * (1.0 - 0.5 * t))  # 255 -> 128
            b = 255
        return max(0, min(255, r)), max(0, min(255, g)), max(0, min(255, b))

    def compute_field(self, x: float, y: float) -> tuple[float, float]:
        Ex = 0.0
        Ey = 0.0
        for c in self.charges:
            dx = x - c.x
            dy = y - c.y
            r2 = dx * dx + dy * dy
            if r2 < R_MIN * R_MIN:
                r2 = R_MIN * R_MIN
            r = math.sqrt(r2)
            inv_r3 = 1.0 / (r2 * r)
            Ex += c.q * dx * inv_r3
            Ey += c.q * dy * inv_r3
        return Ex, Ey

    # =========================================================
    #           ЛИНИИ НАПРЯЖЁННОСТИ
    # =========================================================

    def _integrate_streamline(self, start_x: float, start_y: float) -> list[tuple[float, float]]:
        points: list[tuple[float, float]] = []
        x = start_x
        y = start_y
        end_r2 = FIELDLINE_END_RADIUS * FIELDLINE_END_RADIUS

        for _ in range(STREAMLINE_MAX_STEPS):
            if abs(x) > self.L or abs(y) > self.L:
                break

            for c in self.charges:
                if c.q < 0.0:
                    dx = x - c.x
                    dy = y - c.y
                    if dx * dx + dy * dy < end_r2:
                        points.append((c.x, c.y))
                        return points

            Ex, Ey = self.compute_field(x, y)
            E_mag = math.hypot(Ex, Ey)
            if E_mag < 1e-4:
                break

            Ex /= E_mag
            Ey /= E_mag
            points.append((x, y))
            x += STREAMLINE_STEP * Ex
            y += STREAMLINE_STEP * Ey

        return points

    def rebuild_streamlines(self) -> None:
        for line in self.streamline_shapes:
            line.delete()
        self.streamline_shapes = []

        for c in self.charges:
            if c.q <= 0.0:
                continue

            num_lines = max(20, int(20 * abs(c.q)))
            base_radius = R_MIN * 1.5

            for k in range(num_lines):
                angle = 2.0 * math.pi * k / num_lines
                start_x = c.x + base_radius * math.cos(angle)
                start_y = c.y + base_radius * math.sin(angle)

                world_points = self._integrate_streamline(start_x, start_y)
                if len(world_points) < 2:
                    continue

                screen_points = [
                    self.world_to_screen(px, py) for (px, py) in world_points
                ]

                line = shapes.MultiLine(
                    *screen_points,
                    closed=False,
                    thickness=2.0,
                    color=(0, 0, 0, 255),
                    batch=self.batch,
                    group=self.group_fieldlines,
                )
                self.streamline_shapes.append(line)

    # =========================================================
    #           ИГРОВАЯ ПЕТЛЯ
    # =========================================================

    def update_simulation(self, dt: float) -> None:
        self.time_accumulator += dt
        while self.time_accumulator >= PHYSICS_DT:
            self._physics_step(PHYSICS_DT)
            self.time_accumulator -= PHYSICS_DT

        self.time_since_field_update += dt
        if self.time_since_field_update >= FIELD_UPDATE_INTERVAL:
            self.compute_potential_grid()
            self.rebuild_streamlines()
            self.time_since_field_update = 0.0

        for c in self.charges:
            sx, sy = self.world_to_screen(c.x, c.y)
            c.shape.position = (sx, sy)

    def _physics_step(self, dt: float) -> None:
        if not self.charges:
            return

        forces = self.compute_forces()

        for c, (fx, fy) in zip(self.charges, forces):
            c.apply_force(fx, fy, dt, DAMPING)

        for c in self.charges:
            c.update_position(dt)

        self.merge_charges()

    def on_draw(self) -> None:
        self.clear()
        self.batch.draw()


if __name__ == "__main__":
    window = FieldSimulationWindow()
    pyglet.app.run()
