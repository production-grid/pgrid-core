export default {

  swap(items, leftIndex, rightIndex){
      var temp = items[leftIndex];
      items[leftIndex] = items[rightIndex];
      items[rightIndex] = temp;
  },

  partition(items, left, right) {
    var pivot   = items[Math.floor((right + left) / 2)].index, //middle element
        i       = left, //left pointer
        j       = right; //right pointer
    while (i <= j) {
        while (items[i].index < pivot) {
            i++;
        }
        while (items[j].index > pivot) {
            j--;
        }
        if (i <= j) {
            this.swap(items, i, j); //sawpping two elements
            i++;
            j--;
        }
    }
    return i;
  },

  quickSortAll(items) {

    return this.quickSort(items, 0, items.length - 1)

  },

  quickSort(items, left, right) {
      var index;
      if (items.length > 1) {
          index = this.partition(items, left, right);
          if (left < index - 1) {
              this.quickSort(items, left, index - 1);
          }
          if (index < right) {
              this.quickSort(items, index, right);
          }
      }
      return items;
  }


}
