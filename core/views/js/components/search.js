export function search() {
  return {
    url: '/tv-show?id=',
    ttID: '',
    isLoading: false,
    isError: false,
    submit() {
      this.isLoading = true;
      fetch(this.url+this.ttID).then(response => {
        if (response.ok) {
          window.location.href = this.url+this.ttID;
        } else {
          this.isLoading = false;
          this.isError = true;
        }
      });
    }
  }
}