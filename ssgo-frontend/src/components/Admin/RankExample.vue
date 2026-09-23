<template>
  <div id="exampleRank" class="modal" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">
            輸入全級排名例子
          </h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <textarea
            rows="7"
            class="form-control mb-2"
            v-model="example"
            disabled
          >
          </textarea>
          <div class="alert alert-info" role="alert">
            可以利用 <code>rank 0</code> 表示重讀學生
          </div>
        </div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-secondary"
            data-bs-dismiss="modal"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      example:
        'userAlias,rank\nlp00000001,3\nlp00000002,4\nlp00000003,6\nlp00000004,5\nlp00000005,2\nlp00000006,1'
    }
  }
}
</script>
