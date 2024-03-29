<script>
export default {
    props: ["user_data"],

    data: function () {
        return {
            modal_data: [],

            data_type: "followers",

            // Dynamic loading parameters
            data_ended: false,
            start_idx: 0,
            limit: 10,

            loading: false,
        };
    },
    methods: {
        // Visit the profile of the user with the given id
        visit(user_id) {
            this.$router.push({ path: "/profile/" + user_id })
        },

        // Reset the current data and fetch the first batch of users
        // then show the modal
        async loadData(type) {
            // Reset the parameters and the users array
            this.data_type = type
            this.start_idx = 0
            this.data_ended = false
            this.modal_data = []

            // Fetch the first batch of users
            let status = await this.loadContent()

            // Show the modal if the request was successful
            if (status) this.$refs["mymodal"].showModal()
            // If the request fails, the interceptor will show the error modal
        },

        // Fetch users from the server
        async loadContent() {
            // Fetch followers / following from the server
            // uses /followers and /following endpoints
            let response = await this.$axios.get("/users/" + this.user_data["user_id"] + "/" + this.data_type + "?start_index=" + this.start_idx + "&limit=" + this.limit)
            if (response == null) return false // An error occurred. The interceptor will show a modal

            // If the server returned less elements than requested,
            // it means that there are no more photos to load
            if (response.data.length == 0 || response.data.length < this.limit)
                this.data_ended = true

            // Append the new photos to the array
            this.modal_data = this.modal_data.concat(response.data)
            return true
        },

        // Load more users when the user scrolls to the bottom
        loadMore() {
            // Avoid sending a request if there are no more photos
            if (this.loading || this.data_ended) return

            // Increase the start index and load more photos
            this.start_idx += this.limit
            this.loadContent()
        },
    },
}
</script>

<template>

    <!-- Modal to show the followers / following -->
    <Modal ref="mymodal" id="userModal" :title="data_type">
        <ul>
            <li v-for="item in modal_data" :key="item.user_id" class="mb-2" style="cursor: pointer"
                @click="visit(item.user_id)" data-bs-dismiss="modal">
                <h5>{{ item.name }}</h5>
            </li>
            <IntersectionObserver sentinal-name="load-more-users" @on-intersection-element="loadMore" />
        </ul>
    </Modal>

    <!-- Profile counters -->
    <div class="row text-center mt-2 mb-3">

        <!-- Photos counter -->
        <div class="col-4" style="border-right: 1px">
            <h3>{{ user_data["photos"] }}</h3>
            <h6>Photos</h6>
        </div>

        <!-- Followers counter -->
        <div class="col-4" @click="loadData('followers')" style="cursor: pointer">
            <h3>{{ user_data["followers"] }}</h3>
            <h6>Followers</h6>
        </div>

        <!-- Following counter -->
        <div class="col-4" @click="loadData('following')" style="cursor: pointer">
            <h3>{{ user_data["following"] }}</h3>
            <h6>Following</h6>
        </div>
    </div>

</template>
