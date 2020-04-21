import { Pipe, PipeTransform } from '@angular/core';

const colors = [
  '#000000',	'#FF0000',	'#00FF00',	'#FFFF00',
  '#0000FF',	'#00FFFF',	'#FF00FF',	'#FFFFFF',
  '#FF7F00',	'#7F7F7F',	'#BFBFBF',	'#007F00',
  '#7FFF00',	'#00007F',	'#7F0000',	'#7F4000',
  '#FF9933',	'#007F7F',	'#7F007F',	'#007F7F',
  '#7F00FF',	'#3399CC',	'#CCFFCC',	'#006633',
  '#FF0033',	'#B21919',	'#993300',	'#CC9933',
  '#999933',	'#FFFFBF',	'#FFFF7F'
];

@Pipe({
  name: 'q3color'
})
export class Q3ColorPipe implements PipeTransform {
  transform(value: string, noColor: boolean = false): string {
    let color = -1;
    let nextColor = 7;
    let buffer = '';

    for (let i = 0; i < value.length; i++) {
      let colorChar = -1;
      if (value[i] === '^') {
        if (value[i + 1] !== '^') {
          colorChar = i + 1;
        } else if (value[i + 2] === value[i + 3]) {
          colorChar = i + 2;
        }
      }

      if (colorChar !== -1) {
        // tslint:disable-next-line: no-bitwise
        nextColor = (value.charCodeAt(colorChar) + 16) & 31;
        i = colorChar + 1;
        continue;
      }

      if (color !== nextColor) {
        if (color !== -1 && !noColor) {
          buffer += '</font>';
        }

        color = nextColor;
        if (!noColor) {
          buffer += '<font color="' + colors[color] + '">';
        }
      }

      buffer += value[i].replace(/<.*?>/g, '');
    }

    if (color !== -1 && !noColor) {
      buffer += '</font>';
    }

    return buffer;
  }
}
