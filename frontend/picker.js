function createColorPicker(buttonEl, targetEl, defaultColor, cssProperty = 'color') {
    const pickr = Pickr.create({
        el: buttonEl,
        useAsButton: true,
        theme: 'classic',
        default: defaultColor,
        swatches: [
            'rgba(244, 67, 54, 1)',
            'rgba(233, 30, 99, 0.95)',
            'rgba(156, 39, 176, 0.9)',
            'rgba(103, 58, 183, 0.85)',
            'rgba(63, 81, 181, 0.8)',
            'rgba(33, 150, 243, 0.75)',
            'rgba(3, 169, 244, 0.7)',
            'rgba(0, 188, 212, 0.7)',
            'rgba(0, 150, 136, 0.75)',
            'rgba(76, 175, 80, 0.8)',
            'rgba(139, 195, 74, 0.85)',
            'rgba(205, 220, 57, 0.9)',
            'rgba(255, 235, 59, 0.95)',
            'rgba(255, 193, 7, 1)'
        ],
        components: {
            preview: true,
            opacity: true,
            hue: true,
            interaction: {
             hex: true,
                rgba: true,
                hsla: false,
                hsva: false,
                cmyk: true,
                input: true,
                clear: true,
                save: true
            }
        }
    });

    pickr.on('change', (color) => {
        targetEl.style[cssProperty] = color.toRGBA().toString();
    });
}

/* H1 */
document.querySelectorAll('.topBanner').forEach(div => {
    const h1 = div.querySelector('h1');
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, h1, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.topBanner').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#000000', 'backgroundColor');  
});

document.querySelectorAll('.topBanner').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#FF7F50', 'borderColor');  
});

/* TOTAL SAVINGS WIDGET */

document.querySelectorAll('.totalSavings').forEach(div => {
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, div, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.totalSavings').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#625DFF17', 'backgroundColor');  
});

document.querySelectorAll('.totalSavings').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#625DFF17', 'borderColor');  
});

/* MAIN PIE WIDGET */

document.querySelectorAll('.mainPie').forEach(div => {
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, div, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.mainPie').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#FF005712', 'backgroundColor');  
});

document.querySelectorAll('.mainPie').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#FF005712', 'borderColor');  
});

/* CYCLING PIE WIDGET */

document.querySelectorAll('.cyclePie').forEach(div => {
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, div, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.cyclePie').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#7000841C', 'backgroundColor');  
});

document.querySelectorAll('.cyclePie').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#7000841C', 'borderColor');  
});


/* STOCK WIDGET */

document.querySelectorAll('.stockWidget').forEach(div => {
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, div, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.stockWidget').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#00664838', 'backgroundColor');  
});

document.querySelectorAll('.stockWidget').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#00664838', 'borderColor');  
});


/* BAR CHART WIDGET */

document.querySelectorAll('.monthBar').forEach(div => {
    const fontBtn = div.querySelector('.font-color-btn');
    createColorPicker(fontBtn, div, '#FFFFFF', 'color'); 
});

document.querySelectorAll('.monthBar').forEach(div => {
    const fillBtn = div.querySelector('.fill-btn');
    createColorPicker(fillBtn, div, '#00415E17', 'backgroundColor');  
});

document.querySelectorAll('.monthBar').forEach(div => {
    const fillBtn = div.querySelector('.border-color-btn');
    createColorPicker(fillBtn, div, '#00415E17', 'borderColor');  
});

/* BACKGROUND */

const bodyFillBtn = document.querySelector('.fill-btn');
createColorPicker(bodyFillBtn, document.body, '#0E1319', 'backgroundColor');


