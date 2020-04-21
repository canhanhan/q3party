import { Pipe, PipeTransform } from '@angular/core';

const knownMaps: string[] = ['Q3CTF1', 'Q3CTF2', 'Q3CTF3', 'Q3CTF4', 'Q3DM0', 'Q3DM1', 'Q3DM10', 'Q3DM11', 'Q3DM12', 'Q3DM13', 'Q3DM14', 'Q3DM15', 'Q3DM16', 'Q3DM17', 'Q3DM18', 'Q3DM19', 'Q3DM2', 'Q3DM3', 'Q3DM4', 'Q3DM5', 'Q3DM6', 'Q3DM7', 'Q3DM8', 'Q3DM9', 'Q3TOURNEY1', 'Q3TOURNEY2', 'Q3TOURNEY3', 'Q3TOURNEY4', 'Q3TOURNEY5', 'Q3TOURNEY6'];

@Pipe({
  name: 'mapName'
})
export class MapNamePipe implements PipeTransform {
  transform(value: string): unknown {
    const index = knownMaps.indexOf(value.toUpperCase());
    if (index > -1) {
      return knownMaps[index];
    }

    return 'UNKNOWNMAP';
  }
}
