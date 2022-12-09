<script>

export default {
	props: ["user_id", "name", "followed", "banned", "my_id", "show_new_post"],
	data: function() {
		return {
			errorMsg: "aaa",
			user_followed: this.post_followed,
			user_banned: this.banned,
            myself: this.my_id == this.user_id,
            showModal: false,
		}
	},
	methods: {
        visit() {
            this.$router.push({ path: "/profile/" + this.user_id });
        },
		follow() {
			this.$axios.put("/users/" + this.user_id + "/followers/" + this.my_id)
            .then(response => {
                this.user_followed = true
                this.$emit('updateInfo')
            })
            .catch(error => alert(error.toString()));
		},
		unfollow() {
			this.$axios.delete("/users/" + this.user_id + "/followers/" + this.my_id)
            .then(response => {
                this.user_followed = false
                this.$emit('updateInfo')
            })
            .catch(error => alert(error.toString()));
		},
		ban() {
			this.$axios.put("/users/" + this.my_id + "/bans/" + this.user_id)
            .then(response => {
                this.user_banned = true
                this.$emit('updateInfo')
            })
            .catch(error => alert(error.toString()));
		},
		unban() {
			this.$axios.delete("/users/" + this.my_id + "/bans/" + this.user_id)
            .then(response => {
                this.user_banned = true
                this.$emit('updateInfo')
            })
            .catch(error => alert(error.toString()));
		},
        openModal() {
            var modal = document.getElementById("exampleModal1");
            modal.modal();
        },
	},
    created() {

    },
}
</script>

<template>
    <div class="card mb-3">

		<div class="container">
			<div class="row">
				<div class="col-10">
					<div class="card-body h-100 d-flex align-items-center">
						<a @click="visit"><h5 class="card-title mb-0">{{ name }}</h5></a>
					</div>
				</div>

				<div class="col-2">
					<div class="card-body d-flex justify-content-end">
                        <div v-if="!myself" class="d-flex">
                            <button v-if="!user_banned" @click="ban" type="button" class="btn btn-outline-danger me-2">Ban</button>
                            <button v-if="user_banned" @click="unban" type="button" class="btn btn-danger me-2">Banned</button>
                            <button v-if="!user_followed" @click="follow" type="button" class="btn btn-outline-primary">Follow</button>
                            <button v-if="user_followed" @click="unfollow" type="button" class="btn btn-primary">Following</button>
                        </div>
                        <div v-if="(myself && !show_new_post)">
                            <button disabled type="button" class="btn btn-secondary">Yourself</button>
                        </div>
                        <div v-if="(myself && show_new_post)" class="d-flex">
                            <button type="button" class="btn btn-primary" @click="showModal = true">Post</button>


                            
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>