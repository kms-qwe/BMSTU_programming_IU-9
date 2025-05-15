import glfw
from OpenGL.GL import *
from OpenGL.GLU import *
from OpenGL.GLUT import *
import numpy as np
import math
import time
import random

# --- Optimization modes ---
modes = ["None", "Backface Culling", "Display List", "Mipmapping Control", "Vertex Arrays"]
mode_index = 0

# --- Scene parameters ---
bounds = 2.0
num_objects = 1

# --- FPS tracking ---
fps = 0.0
frame_count = 0
fps_time = time.time()

# --- Display list ID ---
dl_id = None

# --- Vertex array data ---
vertex_array = None
normal_array = None
texcoord_array = None
primitive_count = 0

# --- Pyramid objects ---
class PyramidObj:
    def __init__(self):
        self.pos = np.array([random.uniform(-bounds, bounds) for _ in range(3)], dtype=float)
        self.vel = np.array([random.uniform(-1.5, 1.5) for _ in range(3)], dtype=float)

pyramids = []

# --- Generate twisted pyramid vertices ---
def generate_twisted_pyramid_vertices(n_sides=4, n_slices=20, base_radius=1.0, height=2.0, twist_angle_total=2*math.pi):
    verts = []
    for i in range(n_slices):
        t = i / (n_slices - 1)
        scale = 1 - t
        twist = twist_angle_total * t
        slice_verts = []
        for j in range(n_sides):
            angle = 2*math.pi*j/n_sides + twist
            x = base_radius*scale*math.cos(angle)
            y = base_radius*scale*math.sin(angle)
            z = height * t
            slice_verts.append((x, y, z))
        verts.append(slice_verts)
    return verts

# --- Immediate-mode draw (for display list) ---
def draw_twisted_pyramid_immediate():
    verts = generate_twisted_pyramid_vertices()
    n_sides = len(verts[0])
    n_slices = len(verts)
    # sides
    for i in range(n_slices - 1):
        primitive = GL_TRIANGLES if (i == n_slices - 2) else GL_QUADS
        glBegin(primitive)
        for j in range(n_sides):
            next_j = (j + 1) % n_sides
            v0 = np.array(verts[i][j])
            v1 = np.array(verts[i+1][j])
            v2 = np.array(verts[i][next_j])
            normal = np.cross(v1 - v0, v2 - v0)
            normal = normal / np.linalg.norm(normal)
            glNormal3fv(normal)
            glTexCoord2f(j/(n_sides-1), i/(n_slices-1)); glVertex3fv(verts[i][j])
            if primitive == GL_QUADS:
                glTexCoord2f(j/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][j])
                glTexCoord2f((j+1)/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][next_j])
                glTexCoord2f((j+1)/(n_sides-1), i/(n_slices-1)); glVertex3fv(verts[i][next_j])
            else:
                glTexCoord2f((j+1)/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i][next_j])
                glTexCoord2f(j/(n_sides-1), (i+1)/(n_slices-1)); glVertex3fv(verts[i+1][j])
        glEnd()
    # base
    glBegin(GL_TRIANGLES)
    center = (0.0, 0.0, 0.0)
    for j in range(n_sides):
        next_j = (j + 1) % n_sides
        glNormal3f(0, 0, -1)
        glTexCoord2f(0.5, 0.5); glVertex3fv(center)
        glTexCoord2f(j/(n_sides-1), 0); glVertex3fv(verts[0][j])
        glTexCoord2f((j+1)/(n_sides-1), 0); glVertex3fv(verts[0][next_j])
    glEnd()

# --- Initialize vertex arrays ---
def init_vertex_arrays():
    global vertex_array, normal_array, texcoord_array, primitive_count
    verts = generate_twisted_pyramid_vertices()
    n_sides = len(verts[0])
    n_slices = len(verts)
    data = []
    # sides: convert quads to triangles
    for i in range(n_slices - 1):
        for j in range(n_sides):
            next_j = (j + 1) % n_sides
            # two triangles per quad except last slice
            if i < n_slices - 2:
                tri_pts = [(i,j), (i+1,j), (i+1,next_j), (i,j), (i+1,next_j), (i,next_j)]
            else:
                tri_pts = [(i,j), (i+1,j), (i+1,next_j)]
            for (si, sj) in tri_pts:
                v0 = np.array(verts[si][sj])
                # compute normal of triangle
                # approximate: use face normal for all three
                # but for quads we've split quads; use first three
                # pick next two vertices
                # simple: reuse immediate mode normal
                # compute per-vertex normal
                # here compute per-face
                # get two edges
                # but we skip normal precision; compute for each triangle
                # fetch next two points for normal
                # For triangles, flatten
                # For simplicity, recalc normal for each vertex same
                # implement:
                # determine triangle coords
                # but easier: reuse normal calculation
                # so calculate tri vertices
                if tri_pts.index((si,sj)) == 0:
                    # compute once
                    v1 = np.array(verts[tri_pts[1][0]][tri_pts[1][1]])
                    v2 = np.array(verts[tri_pts[2][0]][tri_pts[2][1]])
                    normal = np.cross(v1 - v0, v2 - v0)
                    normal = normal / np.linalg.norm(normal)
                data.append((*v0, *normal, sj/(n_sides-1), si/(n_slices-1)))
    # convert to numpy
    arr = np.array(data, dtype=np.float32)
    # columns: x,y,z,nx,ny,nz,u,v
    vertex_array = arr[:, 0:3]
    normal_array = arr[:, 3:6]
    texcoord_array = arr[:, 6:8]
    primitive_count = len(vertex_array)

# --- Draw via vertex arrays ---
def draw_twisted_pyramid_va():
    glEnableClientState(GL_VERTEX_ARRAY)
    glEnableClientState(GL_NORMAL_ARRAY)
    glEnableClientState(GL_TEXTURE_COORD_ARRAY)
    glVertexPointer(3, GL_FLOAT, 0, vertex_array)
    glNormalPointer(GL_FLOAT, 0, normal_array)
    glTexCoordPointer(2, GL_FLOAT, 0, texcoord_array)
    glDrawArrays(GL_TRIANGLES, 0, primitive_count)
    glDisableClientState(GL_TEXTURE_COORD_ARRAY)
    glDisableClientState(GL_NORMAL_ARRAY)
    glDisableClientState(GL_VERTEX_ARRAY)

# --- Lighting and material (reuse) ---
light_params = [
    {"id": GL_LIGHT0, "ambient": [0.2,0.0,0.0,1.0], "diffuse": [1.0,0.0,0.0,1.0],
     "specular": [1.0,0.2,0.2,1.0], "pos": [3.0,3.0,3.0,1.0], "color": (1.0,0.0,0.0)},
    {"id": GL_LIGHT1, "ambient": [0.0,0.0,0.2,1.0], "diffuse": [0.0,0.0,1.0,1.0],
     "specular": [0.2,0.2,1.0,1.0], "pos": [-3.0,-3.0,2.0,1.0], "color": (0.0,0.0,1.0)}
]
def init_lighting_and_material():
    glEnable(GL_LIGHTING)
    glLightModelfv(GL_LIGHT_MODEL_AMBIENT, [0.05,0.05,0.05,1.0])
    for lp in light_params:
        glEnable(lp["id"])
        glLightfv(lp["id"], GL_AMBIENT, lp["ambient"])
        glLightfv(lp["id"], GL_DIFFUSE, lp["diffuse"])
        glLightfv(lp["id"], GL_SPECULAR, lp["specular"])
        glLightfv(lp["id"], GL_POSITION, lp["pos"])
    glMaterialfv(GL_FRONT_AND_BACK, GL_AMBIENT, [0.2,0.2,0.2,1.0])
    glMaterialfv(GL_FRONT_AND_BACK, GL_DIFFUSE, [0.8,0.5,0.3,1.0])
    glMaterialfv(GL_FRONT_AND_BACK, GL_SPECULAR, [1.0,1.0,1.0,1.0])
    glMaterialf(GL_FRONT_AND_BACK, GL_SHININESS, 50.0)

# --- Texture (reuse + mipmaps) ---
def load_procedural_texture():
    glEnable(GL_TEXTURE_2D)
    glTexEnvf(GL_TEXTURE_ENV, GL_TEXTURE_ENV_MODE, GL_MODULATE)
    tex_size = 64
    checker = np.zeros((tex_size,tex_size), dtype=np.uint8)
    for i in range(tex_size):
        for j in range(tex_size):
            checker[i,j] = 255 if (i//8 + j//8) % 2 == 0 else 128
    glTexImage2D(GL_TEXTURE_2D, 0, GL_LUMINANCE, tex_size, tex_size, 0,
                 GL_LUMINANCE, GL_UNSIGNED_BYTE, checker)
    glGenerateMipmap(GL_TEXTURE_2D)
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)

# --- Draw bounding box and lights (reuse) ---
def draw_bounding_box(size=2.0, alpha=0.3):
    glDisable(GL_LIGHTING)
    glColor4f(1,1,1,alpha);
    glLineWidth(2.0)
    glBegin(GL_LINES)
    for x in (-size,size):
        for y in (-size,size): glVertex3f(x,y,-size); glVertex3f(x,y,size)
    for x in (-size,size):
        for z in (-size,size): glVertex3f(x,-size,z); glVertex3f(x,size,z)
    for y in (-size,size):
        for z in (-size,size): glVertex3f(-size,y,z); glVertex3f(size,y,z)
    glEnd()
    glEnable(GL_LIGHTING)

def draw_light_sources():
    glDisable(GL_LIGHTING)
    quad = gluNewQuadric()
    for lp in light_params:
        glPushMatrix(); glTranslatef(*lp["pos"][:3])
        glColor3f(*lp["color"]); gluSphere(quad,0.1,16,16)
        glPopMatrix()
    glEnable(GL_LIGHTING)

# --- 2D text overlay ---
def draw_text(window, x, y, text):
    w,h = glfw.get_framebuffer_size(window)
    glMatrixMode(GL_PROJECTION); glPushMatrix(); glLoadIdentity()
    glOrtho(0, w, 0, h, -1, 1)
    glMatrixMode(GL_MODELVIEW); glPushMatrix(); glLoadIdentity()
    glDisable(GL_LIGHTING)
    glColor3f(1,1,1)
    glRasterPos2f(x, y)
    for c in text: glutBitmapCharacter(GLUT_BITMAP_HELVETICA_18, ord(c))
    glEnable(GL_LIGHTING)
    glPopMatrix(); glMatrixMode(GL_PROJECTION); glPopMatrix(); glMatrixMode(GL_MODELVIEW)

# --- GLFW key callback ---
def key_callback(window, key, scancode, action, mods):
    global mode_index, num_objects, pyramids
    if action == glfw.PRESS:
        if key == glfw.KEY_O:
            mode_index = (mode_index + 1) % len(modes)
            # backface culling
            if modes[mode_index] == "Backface Culling": glEnable(GL_CULL_FACE); glCullFace(GL_BACK)
            else: glDisable(GL_CULL_FACE)
            # mipmapping control
            if modes[mode_index] == "Mipmapping Control":
                glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR_MIPMAP_LINEAR)
                glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR)
            else:
                glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
                glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
        elif key in (glfw.KEY_EQUAL, glfw.KEY_KP_ADD):
            num_objects += 1; pyramids = [PyramidObj() for _ in range(num_objects)]
        elif key in (glfw.KEY_MINUS, glfw.KEY_KP_SUBTRACT):
            if num_objects > 1:
                num_objects -= 1; pyramids = [PyramidObj() for _ in range(num_objects)]

# --- Main ---
def main():
    global dl_id, pyramids, fps_time, frame_count, fps
    # init GLUT for text
    glutInit()
    if not glfw.init(): return
    window = glfw.create_window(800,600, "Lab7: OpenGL Optimization", None, None)
    if not window: glfw.terminate(); return
    glfw.make_context_current(window)
    glfw.set_key_callback(window, key_callback)

    glEnable(GL_DEPTH_TEST)
    glMatrixMode(GL_PROJECTION); glLoadIdentity(); gluPerspective(45,800/600,0.1,50)
    glMatrixMode(GL_MODELVIEW)

    init_lighting_and_material()
    load_procedural_texture()
    # create display list
    dl_id = glGenLists(1); glNewList(dl_id, GL_COMPILE); draw_twisted_pyramid_immediate(); glEndList()
    init_vertex_arrays()
    pyramids = [PyramidObj() for _ in range(num_objects)]

    last_t = time.time()
    while not glfw.window_should_close(window):
        t = time.time(); dt = t - last_t; last_t = t
        # FPS
        frame_count += 1
        if t - fps_time >= 1.0:
            fps = frame_count / (t - fps_time)
            fps_time = t; frame_count = 0

        # update pyramids
        for obj in pyramids:
            obj.pos += obj.vel * dt
            for i in range(3):
                if abs(obj.pos[i]) > bounds:
                    obj.pos[i] = np.sign(obj.pos[i]) * bounds; obj.vel[i] = -obj.vel[i]

        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT)
        glLoadIdentity(); gluLookAt(5,5,5, 0,0,0, 0,0,1)
        for lp in light_params: glLightfv(lp["id"], GL_POSITION, lp["pos"])
        draw_bounding_box(bounds, 0.3)
        draw_light_sources()

        # draw each pyramid
        for obj in pyramids:
            glPushMatrix()
            glTranslatef(*obj.pos)
            glRotatef((t*30)%360, 0,0,1)
            mode = modes[mode_index]
            if mode == "Display List": glCallList(dl_id)
            elif mode == "Vertex Arrays": draw_twisted_pyramid_va()
            else: draw_twisted_pyramid_immediate()
            glPopMatrix()

        # overlay text
        draw_text(window, 10, 10, f"Mode: {modes[mode_index]}  Objects: {num_objects}  FPS: {fps:.1f}")

        glfw.swap_buffers(window)
        glfw.poll_events()

    glfw.terminate()

if __name__ == "__main__": main()
