import dayjs from "dayjs";

export function humanizeDate(date: string | Date | null): string {
    if (date === null) {
        return '-'
    }

    const d = dayjs(date)
    if (d.year() <= 1 || d.year() === 1970) {
        return '-'
    }
    return d.format('YYYY-MM-DD HH:mm:ss')
}
