import glfw
from OpenGL.GL import *
import numpy as np
import math
import time
import ctypes

# ----------------------------------------
# GLSL-шейдеры
# ----------------------------------------
VERT_SHADER = """
#version 330 core
layout(location=0) in vec3 in_pos;

layout(location=1) in vec3 in_norm;
layout(location=2) in vec2 in_uv;

uniform mat4 model;
uniform mat4 view;
uniform mat4 proj;
uniform mat3 normal_matrix;

out vec3 fs_pos;
out vec3 fs_norm;
out vec2 fs_uv;

void main() {
    vec4 world_pos = model * vec4(in_pos, 1.0);
    fs_pos = world_pos.xyz;
    fs_norm = normalize(normal_matrix * in_norm);
    fs_uv  = in_uv;
    gl_Position = proj * view * world_pos;
}
"""

FRAG_SHADER = """
#version 330 core
struct Light {
    vec3  pos;
    vec3  ambient;
    vec3  diffuse;
    vec3  specular;
};
uniform Light lights[2];
uniform vec3 view_pos;
uniform vec3 mat_ambient, mat_diffuse, mat_specular;
uniform float mat_shininess;
uniform sampler2D checker_tex;

in vec3 fs_pos;
in vec3 fs_norm;
in vec2 fs_uv;
out vec4 out_color;

void main() {
    // выборка процедурного фактора (luminance)
    float l = texture(checker_tex, fs_uv).r;

    vec3 N = normalize(fs_norm);
    vec3 V = normalize(view_pos - fs_pos);
    vec3 result = vec3(0.0);

    for(int i = 0; i < 2; ++i){
        vec3 L = normalize(lights[i].pos - fs_pos);
        // ambient
        vec3 ambient = lights[i].ambient * mat_ambient;
        // diffuse
        float diff = max(dot(N, L), 0.0);
        vec3 diffuse = diff * lights[i].diffuse * mat_diffuse;
        // specular
        vec3 H = normalize(L + V);
        float spec = pow(max(dot(N, H), 0.0), mat_shininess);
        vec3 specular = spec * lights[i].specular * mat_specular;
        result += ambient + diffuse + specular;
    }

    result *= l; 
    out_color = vec4(result, 1.0);
}
"""

# ----------------------------------------
# Утилиты: компиляция шейдеров и линковка
# ----------------------------------------
def compile_shader(src, kind):
    sh = glCreateShader(kind)
    glShaderSource(sh, src)
    glCompileShader(sh)
    if not glGetShaderiv(sh, GL_COMPILE_STATUS):
        err = glGetShaderInfoLog(sh).decode()
        raise RuntimeError(f"Shader compile error: {err}")
    return sh

def create_program():
    v = compile_shader(VERT_SHADER, GL_VERTEX_SHADER)
    f = compile_shader(FRAG_SHADER, GL_FRAGMENT_SHADER)
    pr = glCreateProgram()
    glAttachShader(pr, v)
    glAttachShader(pr, f)
    glLinkProgram(pr)
    if not glGetProgramiv(pr, GL_LINK_STATUS):
        err = glGetProgramInfoLog(pr).decode()
        raise RuntimeError(f"Program link error: {err}")
    glDeleteShader(v)
    glDeleteShader(f)
    return pr

# ----------------------------------------
# Генерация геометрии закрученной пирамиды
# ----------------------------------------
def generate_twisted_pyramid(n_sides, n_slices, base_r, height, twist_total):
    verts, norms, uvs = [], [], []
    def add_face(v0, v1, v2):
        # рассчитываем нормаль
        n = np.cross(v1-v0, v2-v0)
        n = n / np.linalg.norm(n)
        return n

    # генерируем срезы
    slices = []
    for i in range(n_slices):
        t = i/(n_slices-1)
        r = base_r*(1-t)
        twist = twist_total*t
        ring = []
        for j in range(n_sides):
            ang = 2*math.pi*j/n_sides + twist
            x,y,z = r*math.cos(ang), r*math.sin(ang), height*t
            ring.append((x,y,z))
        slices.append(ring)

    # боковые поверхности
    for i in range(n_slices-1):
        for j in range(n_sides):
            j2 = (j+1)%n_sides
            p0 = np.array(slices[i][j]  )
            p1 = np.array(slices[i+1][j])
            p2 = np.array(slices[i+1][j2])
            p3 = np.array(slices[i][j2] )
            # первая треугольная половина квадрата
            n = add_face(p0,p1,p2)
            for vv,uv in [ (p0,(j/(n_sides-1), i/(n_slices-1))),
                           (p1,(j/(n_sides-1),(i+1)/(n_slices-1))),
                           (p2,((j+1)/(n_sides-1),(i+1)/(n_slices-1))) ]:
                verts.extend(vv); norms.extend(n); uvs.extend(uv)
            # вторая половина
            n = add_face(p2,p3,p0)
            for vv,uv in [ (p2,((j+1)/(n_sides-1),(i+1)/(n_slices-1))),
                           (p3,((j+1)/(n_sides-1), i/(n_slices-1))),
                           (p0,(j/(n_sides-1),           i/(n_slices-1))) ]:
                verts.extend(vv); norms.extend(n); uvs.extend(uv)

    # основание (полигон)
    base = slices[0]
    center = np.array([0,0,0],float)
    up_n = np.array([0,0,-1],float)
    for j in range(1, len(base)-1):
        v0, v1, v2 = base[0], base[j], base[j+1]
        for vv,uv in [ (v0,(0.5,0.5)),
                       (v1,(j/(n_sides-1),0)),
                       (v2,((j+1)/(n_sides-1),0)) ]:
            verts.extend(vv); norms.extend(up_n); uvs.extend(uv)

    return np.array(verts, dtype='f4'), np.array(norms, dtype='f4'), np.array(uvs, dtype='f4')

# ----------------------------------------
# Генерация каркаса куба (линии)
# ----------------------------------------
def generate_box_lines(size):
    s = size
    lines = []
    for x in (-s,s):
        for y in (-s,s):
            lines += [( x,y,-s),( x,y, s)]
    for x in (-s,s):
        for z in (-s,s):
            lines += [( x,-s, z),( x, s, z)]
    for y in (-s,s):
        for z in (-s,s):
            lines += [(-s, y, z),( s, y, z)]
    return np.array(lines, dtype='f4')

# ----------------------------------------
# Генерация UV-сферы для источников
# ----------------------------------------
def generate_sphere(radius, rings, sectors):
    verts = []
    for r in range(rings+1):
        phi = math.pi * r/rings
        for s in range(sectors+1):
            theta = 2*math.pi * s/sectors
            x = radius*math.sin(phi)*math.cos(theta)
            y = radius*math.sin(phi)*math.sin(theta)
            z = radius*math.cos(phi)
            verts.append((x,y,z))
    verts = np.array(verts, dtype='f4')
    # индексы треугольников
    inds = []
    for r in range(rings):
        for s in range(sectors):
            i1 = r*(sectors+1)+s
            i2 = i1 + sectors+1
            inds += [i1,i2,i1+1, i1+1,i2,i2+1]
    return verts, np.array(inds, dtype=np.uint32)

# ----------------------------------------
# Основная функция
# ----------------------------------------
def main():
    if not glfw.init(): return
    glfw.window_hint(glfw.CONTEXT_VERSION_MAJOR,3)
    glfw.window_hint(glfw.CONTEXT_VERSION_MINOR,3)
    glfw.window_hint(glfw.OPENGL_PROFILE,glfw.OPENGL_CORE_PROFILE)
    win = glfw.create_window(800,600,"Lab 6 – shaders",None,None)
    if not win:
        glfw.terminate(); return
    glfw.make_context_current(win)
    glEnable(GL_DEPTH_TEST)

    # создаём шейдерную программу
    prog = create_program()
    glUseProgram(prog)

    # генерируем геометрию пирамиды
    verts, norms, uvs = generate_twisted_pyramid(4,20,1.0,2.0,2*math.pi)
    n_verts = len(verts)//3
    vao_pyr = glGenVertexArrays(1); glBindVertexArray(vao_pyr)
    # VBO общие
    vbo = glGenBuffers(3)
    # позиции
    glBindBuffer(GL_ARRAY_BUFFER, vbo[0])
    glBufferData(GL_ARRAY_BUFFER, verts.nbytes, verts, GL_STATIC_DRAW)
    glEnableVertexAttribArray(0); glVertexAttribPointer(0,3,GL_FLOAT,False,0,None)
    # нормали
    glBindBuffer(GL_ARRAY_BUFFER, vbo[1])
    glBufferData(GL_ARRAY_BUFFER, norms.nbytes, norms, GL_STATIC_DRAW)
    glEnableVertexAttribArray(1); glVertexAttribPointer(1,3,GL_FLOAT,False,0,None)
    # UV
    glBindBuffer(GL_ARRAY_BUFFER, vbo[2])
    glBufferData(GL_ARRAY_BUFFER, uvs.nbytes, uvs, GL_STATIC_DRAW)
    glEnableVertexAttribArray(2); glVertexAttribPointer(2,2,GL_FLOAT,False,0,None)

    # каркас куба
    box = generate_box_lines(2.0)
    vao_box = glGenVertexArrays(1); glBindVertexArray(vao_box)
    vbo_box = glGenBuffers(1)
    glBindBuffer(GL_ARRAY_BUFFER, vbo_box)
    glBufferData(GL_ARRAY_BUFFER, box.nbytes, box, GL_STATIC_DRAW)
    glEnableVertexAttribArray(0); glVertexAttribPointer(0,3,GL_FLOAT,False,0,None)

    # сфера источника
    sph_verts, sph_inds = generate_sphere(0.1,16,16)
    vao_sph = glGenVertexArrays(1); glBindVertexArray(vao_sph)
    vbo_s = glGenBuffers(1); glBindBuffer(GL_ARRAY_BUFFER, vbo_s)
    glBufferData(GL_ARRAY_BUFFER, sph_verts.nbytes, sph_verts, GL_STATIC_DRAW)
    glEnableVertexAttribArray(0); glVertexAttribPointer(0,3,GL_FLOAT,False,0,None)
    ebo_s = glGenBuffers(1); glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, ebo_s)
    glBufferData(GL_ELEMENT_ARRAY_BUFFER, sph_inds.nbytes, sph_inds, GL_STATIC_DRAW)
    n_sph = len(sph_inds)

    # процедурная «шахматка»
    tex = glGenTextures(1)
    glBindTexture(GL_TEXTURE_2D, tex)
    checker = np.zeros((64,64), dtype=np.uint8)
    for i in range(64):
        for j in range(64):
            checker[i,j] = 255 if (i//8+j//8)%2==0 else 128
    glTexImage2D(GL_TEXTURE_2D,0,GL_RED,64,64,0,GL_RED,GL_UNSIGNED_BYTE,checker)
    glTexParameteri(GL_TEXTURE_2D,GL_TEXTURE_MIN_FILTER,GL_NEAREST)
    glTexParameteri(GL_TEXTURE_2D,GL_TEXTURE_MAG_FILTER,GL_NEAREST)
    # влияем на swizzle (чтобы .r → одинаковые rgb)
    swz = np.array([GL_RED,GL_RED,GL_RED,GL_ONE],dtype=np.int32)
    glTexParameterIiv(GL_TEXTURE_2D,GL_TEXTURE_SWIZZLE_RGBA, swz)

    # передаём константы материала
    loc = lambda name: glGetUniformLocation(prog, name)
    glUniform3f(loc("mat_ambient"),  0.2,0.2,0.2)
    glUniform3f(loc("mat_diffuse"),  0.8,0.5,0.3)
    glUniform3f(loc("mat_specular"), 1.0,1.0,1.0)
    glUniform1f(loc("mat_shininess"), 50.0)

    # параметры света
    lights = [
      { "pos":(3,3,3), "ambient":(0.2,0,0), "diffuse":(1,0,0), "specular":(1,0.2,0.2) },
      { "pos":(-3,-3,2),"ambient":(0,0,0.2),"diffuse":(0,0,1),"specular":(0.2,0.2,1) }
    ]
    for i,L in enumerate(lights):
        glUniform3f(loc(f"lights[{i}].pos"),      *L["pos"])
        glUniform3f(loc(f"lights[{i}].ambient"),  *L["ambient"])
        glUniform3f(loc(f"lights[{i}].diffuse"),  *L["diffuse"])
        glUniform3f(loc(f"lights[{i}].specular"), *L["specular"])

    # активируем единицу текстуры
    glUniform1i(loc("checker_tex"), 0)

    last_t = time.time()
    pos = np.array([0,0,0],float)
    vel = np.array([1.2,0.9,1.5],float)
    bounds = 2.0

    while not glfw.window_should_close(win):
        t = time.time(); dt = t - last_t; last_t = t
        pos += vel * dt
        for i in range(3):
            if abs(pos[i])>bounds:
                pos[i] = np.sign(pos[i])*bounds
                vel[i] = -vel[i]

        glClear(GL_COLOR_BUFFER_BIT|GL_DEPTH_BUFFER_BIT)
        # матрицы
        proj = np.array(glfw.get_framebuffer_size(win),dtype=float)
        proj = glm_perspective(45, proj[0]/proj[1],0.1,50.0)
        view = glm_lookat(np.array([5,5,5],float), np.zeros(3), np.array([0,0,1]))
        glUniformMatrix4fv(loc("view"),1,GL_FALSE,view)
        glUniformMatrix4fv(loc("proj"),1,GL_FALSE,proj)
        glUniform3f(loc("view_pos"),5,5,5)

        # 1) каркас куба
        glBindVertexArray(vao_box)
        glUniformMatrix4fv(loc("model"),1,GL_FALSE,np.eye(4))
        inv_nf = np.linalg.inv(np.eye(4)[:3,:3]).T
        glUniformMatrix3fv(loc("normal_matrix"),1,GL_FALSE,inv_nf)
        glDrawArrays(GL_LINES, 0, len(box))

        # 2) источники
        glBindVertexArray(vao_sph)
        for L in lights:
            M = np.eye(4)
            M[:3,3] = L["pos"]
            glUniformMatrix4fv(loc("model"),1,GL_FALSE,M)
            inv_nf = np.linalg.inv(M[:3,:3]).T
            glUniformMatrix3fv(loc("normal_matrix"),1,GL_FALSE,inv_nf)
            glDrawElements(GL_TRIANGLES, n_sph, GL_UNSIGNED_INT, None)

        # 3) вращающаяся закрученная пирамида
        M = np.eye(4)
        # сдвиг
        M[:3,3] = pos
        # вращение вокруг Z
        ang = (t*30)%360 * math.pi/180
        Rz = np.array([[ math.cos(ang), -math.sin(ang),0],
                       [ math.sin(ang),  math.cos(ang),0],
                       [           0,             0,   1]],float)
        M[:3,:3] = Rz @ M[:3,:3]
        glBindVertexArray(vao_pyr)
        glUniformMatrix4fv(loc("model"),1,GL_FALSE,M)
        inv_nf = np.linalg.inv(M[:3,:3]).T
        glUniformMatrix3fv(loc("normal_matrix"),1,GL_FALSE,inv_nf)
        glActiveTexture(GL_TEXTURE0); glBindTexture(GL_TEXTURE_2D, tex)
        glDrawArrays(GL_TRIANGLES, 0, n_verts)

        glfw.swap_buffers(win)
        glfw.poll_events()

    glfw.terminate()

# вспомогательные функции матриц (реализуйте или возьмите из glm-подобного модуля)
def glm_perspective(fovy, aspect, near, far):
    f = 1/math.tan(math.radians(fovy)/2)
    M = np.zeros((4,4),dtype=float)
    M[0,0] = f/aspect
    M[1,1] = f
    M[2,2] = (far+near)/(near-far)
    M[2,3] = (2*far*near)/(near-far)
    M[3,2] = -1
    return M

def glm_lookat(eye, center, up):
    f = center - eye; f /= np.linalg.norm(f)
    u = up/np.linalg.norm(up)
    s = np.cross(f,u); s /= np.linalg.norm(s)
    u = np.cross(s,f)
    M = np.eye(4,dtype=float)
    M[0,:3] = s
    M[1,:3] = u
    M[2,:3] = -f
    M[:3,3] = -M[:3,:3] @ eye
    return M

if __name__=="__main__":
    main()
