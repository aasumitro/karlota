import dayjs from "dayjs";
import isToday from "dayjs/plugin/isToday";
import isYesterday from "dayjs/plugin/isYesterday";

dayjs.extend(isToday);
dayjs.extend(isYesterday);

export const formatOnlineTime = (time?: number) => {
  if (!time) return;
  const date = dayjs.unix(time);
  if (date.isToday()) {
    return "Last online " + date.format("HH:mm");
  } else if (date.isYesterday()) {
    return "Last online yesterday";
  } else if (dayjs().diff(date, "day") <= 5) {
    return "Last online "+ "a few days ago";
  } else {
    return "Last online " + date.format("DD MMM");
  }
};

