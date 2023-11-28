import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'elapsedTime'
})
export class ElapsedTimePipe implements PipeTransform {

  transform(dt: Date | null, now: Date): string {
    if (!dt) {
      return ''
    }
    const elapsed = now.valueOf() - dt.valueOf()
    const minutes = Math.floor(elapsed / (60 * 1000))

    if (minutes > 0) {
      return `${minutes} минут(ы) назад`
    }

    const seconds = Math.floor(elapsed / 1000)
    return `${seconds} секунд(ы) назад`
  }

}
