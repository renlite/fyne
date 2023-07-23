#version 110

/* scaled */
uniform vec2 frame_size;

varying float type;
varying float stroke_width;
varying float stroke_width_half;
varying float radius;

varying vec4 fill_color;
varying vec4 stroke_color;

varying vec4 rect_coords; //x1 [0], x2 [1], y1 [2], y2 [3]; coords of the rect_frame
varying vec2 rect_size_half;

float calc_distance(vec2 p, vec2 b, float r)
{
    vec2 d = abs(p) - b + vec2(r);
	return min(max(d.x, d.y), 0.0) + length(max(d, 0.0)) - r;   
}

void main() {

    if (type == 1.0) {
        vec4 color = fill_color;
        if (gl_FragCoord.x >= rect_coords[1] - stroke_width ){
            color = stroke_color;
        } else if (gl_FragCoord.x <= rect_coords[0] + stroke_width){
            color = stroke_color;
        } else if (gl_FragCoord.y <= frame_size.y - rect_coords[3] + stroke_width ){
            color = stroke_color;
        } else if (gl_FragCoord.y >= frame_size.y - rect_coords[2] - stroke_width ){
            color = stroke_color;
        }
        gl_FragColor = color; 
          
    } else if (type == 2.0 ) {
        vec4 frag_rect_coords = vec4(rect_coords[0], rect_coords[1], frame_size.y - rect_coords[3], frame_size.y - rect_coords[2]);
        vec2 vec_centered_pos = (gl_FragCoord.xy - vec2(frag_rect_coords[0] + frag_rect_coords[1], frag_rect_coords[2] + frag_rect_coords[3]) * 0.5);

        float distance = calc_distance(vec_centered_pos, rect_size_half, radius - stroke_width_half);

        vec4 from_color = stroke_color; //Always the border color. If no border, this still should be set
        vec4 to_color = stroke_color; //Outside color

        if (stroke_width_half == 0.0)
        {
            from_color = fill_color;
            to_color = fill_color;
        }
        to_color[3] = 0.0; // blend the fill colour to alpha

        if (distance < 0.0)
        {
            to_color = fill_color;
        } 

        distance = abs(distance) - stroke_width_half;

        float blend_amount = smoothstep(-1.0, 1.0, distance);

        // final color
        gl_FragColor = mix(from_color, to_color, blend_amount);
    }
}
