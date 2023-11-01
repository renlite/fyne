#version 100

#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
#endif
precision mediump int;
precision lowp sampler2D;
#endif

/* scaled */
// uniform vec2 frame_size; //window_size not used here

// for Texture
uniform sampler2D textures[16];
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
    //not in GL_ES: gl_FragColor = texture2D(textures[idx], fragTexCoord);
    if (idx == 0){
        gl_FragColor = texture2D(textures[0], fragTexCoord);
    } else if (idx == 1){
            gl_FragColor = texture2D(textures[1], fragTexCoord);
    } else if (idx == 2){
            gl_FragColor = texture2D(textures[2], fragTexCoord);
    } else if (idx == 3){
            gl_FragColor = texture2D(textures[3], fragTexCoord);
    } else if (idx == 4){
            gl_FragColor = texture2D(textures[4], fragTexCoord);
    } else if (idx == 5){
            gl_FragColor = texture2D(textures[5], fragTexCoord);
    } else if (idx == 6){
            gl_FragColor = texture2D(textures[6], fragTexCoord);
    } else if (idx == 7){
            gl_FragColor = texture2D(textures[7], fragTexCoord);
    } else if (idx == 8){
            gl_FragColor = texture2D(textures[8], fragTexCoord);
    } else if (idx == 9){
            gl_FragColor = texture2D(textures[9], fragTexCoord);
    } else if (idx == 10){
            gl_FragColor = texture2D(textures[10], fragTexCoord);
    } else if (idx == 11){
            gl_FragColor = texture2D(textures[11], fragTexCoord);
    } else if (idx == 12){
            gl_FragColor = texture2D(textures[12], fragTexCoord);
    } else if (idx == 13){
            gl_FragColor = texture2D(textures[13], fragTexCoord);
    } else if (idx == 14){
            gl_FragColor = texture2D(textures[14], fragTexCoord);
    } else if (idx == 15){
            gl_FragColor = texture2D(textures[15], fragTexCoord);
    }

}
