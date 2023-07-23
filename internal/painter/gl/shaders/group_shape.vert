#version 110

/* scaled */
uniform vec2 frame_size; //window_size

/* triangle 6 points (x/y) for rect */
attribute vec2 att_vert;

/* 
scaled size, coords 
Description of the new shape.
*/
attribute vec4 att_type; // (shape = 1.0, 2.0, 3.0, ... | stroke_width | radius | NOT_USED )
/*
Shapes:
1.0 = rectangle
2.0 = round_rectangle
3.0 = text
*/

attribute vec4 att_fill_color; // (fillColor RGBA)
attribute vec4 att_stroke_color; // (fillColor RGBA)
attribute vec4 att_rect_coords; //x1 [0], x2 [1], y1 [2], y2 [3]; coords of the rect_frame

varying float type;
varying float stroke_width;
varying float stroke_width_half;
varying float radius;

varying vec4  fill_color;
varying vec4  stroke_color;

varying vec4  rect_coords;
varying vec2  rect_size_half;

/*
vec4 unpack_to_frag_color(float rgb_as_float, float a){
    vec4 color;
    color.b = floor(rgb_as_float / 256.0 / 256.0);
    color.g = floor((rgb_as_float - color.r * 256.0 * 256.0) / 256.0);
    color.r = floor(rgb_as_float - color.r * 256.0 * 256.0 - color.g * 256.0);
    color.a = a;
    return color / 256.0;
}
*/

vec2 get_rect_size_half(vec4 rect_coords, float stroke_width){
    vec2 rect_size_half;
    float rect_size_width  = rect_coords[1] - rect_coords[0] - stroke_width;
    float rect_size_height = rect_coords[3] - rect_coords[2] - stroke_width;
    rect_size_half = vec2(rect_size_width*0.5, rect_size_height*0.5);
    return rect_size_half;
}


void main() {
    // normalize att_vert in GPU
    gl_Position = vec4(-1.0 + att_vert.x*2.0/frame_size.x, 1.0 - att_vert.y*2.0/frame_size.y, 0, 1);

    type = att_type[0];
    stroke_width = att_type[1];
    stroke_width_half = stroke_width * 0.5;
    radius = att_type[2];
    fill_color = att_fill_color;
    if (stroke_width != 0.0){       
        stroke_color = att_stroke_color;
    }
    rect_coords = att_rect_coords;
    if (type == 2.0){
        rect_size_half = get_rect_size_half(rect_coords, stroke_width);
    }
}


