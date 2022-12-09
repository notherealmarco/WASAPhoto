<script>

export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked", "my_id"],
	data: function() {
		return {
			imageSrc: "",
			errorMsg: null,
			post_liked: this.liked,
			post_like_cnt: this.likes,
		}
	},
	methods: {
		like() {
			this.$axios.put("/users/" + this.user_id + "/photos/" + this.photo_id + "/likes/" + this.my_id).then(response => {
				this.post_liked = true;
				this.post_like_cnt++;
			}).catch(error => {
				console.log(error);
				this.errorMsg = error.toString();
			});
		},
		unlike() {
			this.$axios.delete("/users/" + this.user_id + "/photos/" + this.photo_id + "/likes/" + this.my_id).then(response => {
				this.post_liked = false;
				this.post_like_cnt--;
			}).catch(error => {
				console.log(error);
				this.errorMsg = error.toString();
			});
		},
	},

	created() {
			this.$axios.get("/users/" + this.user_id + "/photos/" + this.photo_id, {
				responseType: 'arraybuffer'
			}).then(response => {
				const img = document.createElement('img');
				img.src = URL.createObjectURL(new Blob([response.data]));
				img.classList.add("card-img-top");
				this.$refs.imageContainer.appendChild(img);
			});
		},
	}
</script>

<template>
    <div class="card mb-5">
		<!--<img v-auth-img="imageSrc" class="card-img-top" alt="Chicago Skyscrapers"/>-->
		<div ref="imageContainer"></div>

		<div class="container">
			<div class="row">
				<div class="col-10">
					<div class="card-body">
						<h5 class="card-title">{{ name }}</h5>
						<p class="card-text">{{ new Date(Date.parse(date)) }}</p>
					</div>
				</div>

				<div class="col-2">
					<div class="card-body d-flex justify-content-end" style="display: inline-flex"> <!-- not quite sure flex is the right property, but it works -->
						<a><h5><i class="card-title bi bi-chat-right pe-1"></i></h5></a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ comments }}</h6>
						<a v-if="!post_liked" @click="like"><h5><i class="card-title bi bi-suit-heart ps-2 pe-1 like-icon"></i></h5></a>
						<a v-if="post_liked" @click="unlike"><h5><i class="card-title bi bi-heart-fill ps-2 pe-1 like-icon like-red"></i></h5></a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_like_cnt }}</h6>
						<h5></h5>
					</div>
				</div>
			</div>
		</div>

		<!--<ul class="list-group list-group-light list-group-small">
		<li class="list-group-item px-4">Cras justo odio</li>
		<li class="list-group-item px-4">Dapibus ac facilisis in</li>
		<li class="list-group-item px-4">Vestibulum at eros</li>
		</ul>-->
	</div>
	<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
</template>

<style>
.like-icon:hover {
	color: #ff0000;
}
.like-red {
	color: #ff0000;
}
</style>