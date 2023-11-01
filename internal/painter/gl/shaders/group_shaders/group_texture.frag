#version 110

/* scaled */
// uniform vec2 frame_size; //window_size not used here

// for Texture
uniform sampler2D textures[32];
/*
uniform sampler2D texture0;
uniform sampler2D texture1;
uniform sampler2D texture2;
uniform sampler2D texture3;
*/
varying vec2 fragTexCoord;
varying float texIdx;


void main() {
    int idx = int(texIdx);
    gl_FragColor = texture2D(textures[idx], fragTexCoord);

    /* debug
    if (idx == 0){
        //gl_FragColor = vec4(1.0, 0.5, 0.0, 1.0);
        gl_FragColor = texture2D(texture0, fragTexCoord);
    } else if (idx == 1){
        //gl_FragColor = vec4(0.0, 1.0, 0.0, 1.0);
         gl_FragColor = texture2D(texture1, fragTexCoord);
    } else if (idx == 2){
        //gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
         gl_FragColor = texture2D(texture2, fragTexCoord);
    } else if (idx == 3){
        //gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
         gl_FragColor = texture2D(texture3, fragTexCoord);
    }
    */

}
