<script>

export default {
	props: ["user_id", "name", "followed", "banned", "my_id", "show_new_post"],
    watch: { 
        banned: function(new_val, old_val) {
          this.user_banned = new_val;
        },
        followed: function(new_val, old_val) {
          this.user_followed = new_val;
        },
    },
	data: function() {
		return {
			errorMsg: "aaa",
			user_followed: this.followed,
			user_banned: this.banned,
            myself: this.my_id == this.user_id,
            show_post_form: false,
            upload_file: null,
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
                this.user_banned = false
                this.$emit('updateInfo')
            })
            .catch(error => alert(error.toString()));
		},
        load_file(e) {
            let files = e.target.files || e.dataTransfer.files;
            if (!files.length) return;
            this.upload_file = files[0];
        },
        submit_file() {
            this.$axios.post("/users/" + this.my_id + "/photos", this.upload_file)
            .then(response => {
                this.show_post_form = false
                this.$emit('updatePosts')
            })
            .catch(error => alert(error.toString()));
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
                            <button v-if="!user_followed" @click="follow" type="button" class="btn btn-primary">Follow</button>
                            <button v-if="user_followed" @click="unfollow" type="button" class="btn btn-outline-primary">Following</button>
                        </div>
                        <div v-if="(myself && !show_new_post)">
                            <button disabled type="button" class="btn btn-secondary">Yourself</button>
                        </div>
                        <div v-if="(myself && show_new_post)" class="d-flex">
                            <button v-if="!show_post_form" type="button" class="btn btn-primary" @click="show_post_form = true">Post</button>                            
                        </div>
                    </div>
                </div>
            </div>
            <div class="row" v-if="show_post_form">
				<div class="col-9">
					<div class="card-body h-100 d-flex align-items-center">
                        <input @change="load_file" class="form-control form-control-lg" id="formFileLg" type="file" />
					</div>
				</div>

				<div class="col-3">
					<div class="card-body d-flex justify-content-end">
                        <button type="button" class="btn btn-primary btn-lg" @click="submit_file">Publish</button>                            
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>