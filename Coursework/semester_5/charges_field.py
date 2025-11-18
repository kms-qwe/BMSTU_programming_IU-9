import numpy as np
import matplotlib.pyplot as plt


def compute_field_and_potential(charges, xlim, ylim, grid_points, r_min):
    """Вычисляет компоненты поля и потенциал на равномерной сетке."""
    x_vals = np.linspace(xlim[0], xlim[1], grid_points)
    y_vals = np.linspace(ylim[0], ylim[1], grid_points)
    X, Y = np.meshgrid(x_vals, y_vals)
    Ex = np.zeros_like(X)
    Ey = np.zeros_like(Y)
    V = np.zeros_like(X)
    for x_i, y_i, q_i in charges:
        dx = X - x_i
        dy = Y - y_i
        r = np.sqrt(dx * dx + dy * dy)
        r = np.where(r < r_min, r_min, r)
        Ex += q_i * dx / (r ** 3)
        Ey += q_i * dy / (r ** 3)
        V += q_i / r
    return X, Y, Ex, Ey, V


def animate_motion_live(
    charges,
    xlim,
    ylim,
    grid_points,
    r_min,
    dt,
    frames,
    damping,
    merge_threshold=None,
    substeps=1,
    frame_delay=None,
):
    """Анимирует движение зарядов во времени с учётом притяжения и отталкивания."""
    charges_state = [list(c) for c in charges]
    velocities = [[0.0, 0.0] for _ in charges_state]
    if merge_threshold is None:
        merge_threshold = r_min
    if frame_delay is None:
        frame_delay = dt
    x_vals = np.linspace(xlim[0], xlim[1], grid_points)
    y_vals = np.linspace(ylim[0], ylim[1], grid_points)
    X_grid, Y_grid = np.meshgrid(x_vals, y_vals)
    plt.ion()
    fig, ax = plt.subplots(figsize=(8, 8))

    fig.canvas.manager.set_window_title("Визуализация электрических полей")

    cbar = None
    global_max_abs = 0.0
    for frame_idx in range(frames):
        if not plt.fignum_exists(fig.number):
            break
        for _ in range(substeps):
            n = len(charges_state)
            forces = [[0.0, 0.0] for _ in range(n)]
            for i in range(n):
                x_i, y_i, q_i = charges_state[i]
                for j in range(n):
                    if i == j:
                        continue
                    x_j, y_j, q_j = charges_state[j]
                    dx = x_i - x_j
                    dy = y_i - y_j
                    r2 = dx * dx + dy * dy
                    r = r2 ** 0.5
                    if r < r_min:
                        r = r_min
                        r2 = r * r
                    force_mag = q_i * q_j / (r2)
                    fx = force_mag * dx / r
                    fy = force_mag * dy / r
                    forces[i][0] += fx
                    forces[i][1] += fy
            for i in range(n):
                ax_i, ay_i = forces[i]
                vx, vy = velocities[i]
                vx += ax_i * dt
                vy += ay_i * dt
                vx *= damping
                vy *= damping
                velocities[i] = [vx, vy]
                charges_state[i][0] += vx * dt
                charges_state[i][1] += vy * dt
            merge_occurred = True
            while merge_occurred:
                merge_occurred = False
                i = 0
                while i < len(charges_state):
                    j = i + 1
                    while j < len(charges_state):
                        qi = charges_state[i][2]
                        qj = charges_state[j][2]
                        if qi * qj < 0:
                            dx = charges_state[i][0] - charges_state[j][0]
                            dy = charges_state[i][1] - charges_state[j][1]
                            dist = (dx * dx + dy * dy) ** 0.5
                            if dist < merge_threshold:
                                q_new = qi + qj
                                x_new = (charges_state[i][0] + charges_state[j][0]) / 2.0
                                y_new = (charges_state[i][1] + charges_state[j][1]) / 2.0
                                vx_new = (velocities[i][0] + velocities[j][0]) / 2.0
                                vy_new = (velocities[i][1] + velocities[j][1]) / 2.0
                                charges_state.pop(j)
                                velocities.pop(j)
                                charges_state.pop(i)
                                velocities.pop(i)
                                if q_new != 0.0:
                                    charges_state.insert(i, [x_new, y_new, q_new])
                                    velocities.insert(i, [vx_new, vy_new])
                                merge_occurred = True
                                i = -1
                                break
                        j += 1
                    i += 1
        X = X_grid
        Y = Y_grid
        Ex = np.zeros_like(X)
        Ey = np.zeros_like(Y)
        V = np.zeros_like(X)
        for x_i, y_i, q_i in charges_state:
            dx = X - x_i
            dy = Y - y_i
            r = np.sqrt(dx * dx + dy * dy)
            r = np.where(r < r_min, r_min, r)
            Ex += q_i * dx / (r ** 3)
            Ey += q_i * dy / (r ** 3)
            V += q_i / r
        local_max = np.max(np.abs(V))
        if local_max > global_max_abs:
            global_max_abs = local_max
        ax.cla()
        levels = np.linspace(-global_max_abs, global_max_abs, 100)
        cf = ax.contourf(
            X,
            Y,
            V,
            levels=levels,
            cmap="seismic",
            vmin=-global_max_abs,
            vmax=global_max_abs,
            alpha=0.8,
        )
        if cbar is None:
            cbar = fig.colorbar(cf, ax=ax, shrink=0.8)
            cbar.set_label("Potential")
        else:
            cbar.update_normal(cf)
        ax.streamplot(
            X,
            Y,
            Ex,
            Ey,
            color="k",
            linewidth=0.5,
            density=1.2,
            arrowstyle="->",
            arrowsize=1.0,
        )
        for x_i, y_i, q_i in charges_state:
            colour = "red" if q_i > 0 else "blue"
            ax.scatter(x_i, y_i, color=colour, s=80, edgecolors="k", zorder=3)
            ax.text(
                x_i,
                y_i,
                f"{q_i:+}",
                color="white",
                fontsize=9,
                ha="center",
                va="center",
                zorder=4,
            )
        ax.set_xlim(xlim)
        ax.set_ylim(ylim)
        ax.set_aspect("equal")
        ax.set_xlabel("x")
        ax.set_ylabel("y")
        ax.set_title(f"t = {frame_idx * dt * substeps:.2f}")
        plt.tight_layout()
        if plt.fignum_exists(fig.number):
            fig.canvas.draw()
            plt.pause(frame_delay)
        else:
            break
    plt.ioff()


def main():
    """Настройка параметров и запуск анимации."""
    L = 3.0
    xlim = (-L, L)
    ylim = (-L, L)
    grid_points = 200
    r_min = 0.05
    charges = [
        (-1.5, 0.0, +2.0),
        (1.5, 0.0, -2.0),
        (0.0, 1.5, +1.0),
        (0.0, -1.5, -1.0),
    ]
    dt = 0.02
    frames = 120
    damping = 0.95
    merge_threshold = 0.2
    substeps = 5
    frame_delay = 0.02
    animate_motion_live(
        charges,
        xlim,
        ylim,
        grid_points,
        r_min,
        dt,
        frames,
        damping,
        merge_threshold,
        substeps,
        frame_delay,
    )


if __name__ == "__main__":
    main()
