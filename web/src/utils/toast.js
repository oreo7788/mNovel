export function showToast(message) {
  const el = document.createElement('div')
  el.className = 'upload-toast'
  el.textContent = message
  el.setAttribute('role', 'status')
  document.body.appendChild(el)
  setTimeout(() => {
    if (el.parentNode) el.parentNode.removeChild(el)
  }, 2500)
}
