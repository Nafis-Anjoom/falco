export function formatDate(date: Date): string {
  date = new Date();
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0"); // Months are 0-indexed
  const day = date.getDate().toString().padStart(2, "0"); // Add leading zero if necessary

  return `${year}-${month}-${day}`;
}