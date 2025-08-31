import dayjs from "dayjs";

export function humanizeDate(date: string | Date): string {
    const d = dayjs(date)
    if (d.year() <= 1 || d.year() === 1970) {
        return '-'
    }
    return d.format('YYYY-MM-DD HH:mm:ss')
}
