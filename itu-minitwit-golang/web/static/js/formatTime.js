window.addEventListener('DOMContentLoaded', (event) => formatTimeLocal(), false);
function formatTimeLocal() {
  const times = document.getElementsByClassName('time');
  for (const element of times) {
    const datetime = new Date(element.innerHTML);
    const userLocale = navigator.language || navigator.userLanguage;
    const userTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    element.innerHTML = `Written @ ${datetime.toLocaleTimeString(userLocale, {
      timeZone: userTimezone,
    })}, ${datetime.toLocaleDateString(userLocale, {
      timeZone: userTimezone,
    })}`;
  }
}
