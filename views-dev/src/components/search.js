export function search() {
  return {
    url: '/tv-show?id=',
    ttID: '',
    isLoading: false,
    isError: false,
    submit() {
      this.isLoading = true;
      fetch(this.url+this.ttID)
    }
  }
}