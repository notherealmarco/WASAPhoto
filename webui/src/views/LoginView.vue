<script>
import LoadingSpinner from '../components/LoadingSpinner.vue';

export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            some_data: null,
            field_username: "",
        };
    },
    methods: {
        async login() {
            this.loading = true;
            this.errormsg = null;
            try {
                let response = await this.$axios.post("/session", {
                    name: this.field_username,
                });
                //this.$router.push({ name: "home" });
                if (response.status == 201 || response.status == 200) {
                    // Save the token in the session storage
                    sessionStorage.setItem("token", response.data["user_id"]);

                    // Update the header
                    this.$axiosUpdate();

                    this.$router.push({ path: "/" });
                }
                else {
                    this.errormsg = response.data["error"];
                }
            }
            catch (e) {
                this.errormsg = e.toString();
            }
            this.loading = false;
        },
        async refresh() {
            //this.loading = true;
            //this.errormsg = null;
            //try {
            //	let response = await this.$axios.get("/");	
            //	this.some_data = response.data;
            //} catch (e) {
            //	this.errormsg = e.toString();
            //}
            this.loading = false;
        },
    },
    mounted() {
        this.refresh();
    },
    components: { LoadingSpinner }
}
</script>

<template>
	<div class="vh-100 container py-5 h-100">
    <div class="row d-flex justify-content-center align-items-center h-100">
        <!--<div class="col-sm"><h2>* immagina un logo carino *</h2></div>-->
        <div class="col-12 col-md-8 col-lg-6 col-xl-5">
        <div class="card" style="border-radius: 1rem">
        <div class="card-body p-4">

        <h1 class="h2 pb-4 text-center">WASAPhoto</h1>

		<form>
			<!-- Email input -->
			<div class="form-floating mb-4">
				<input v-model="field_username" type="email" id="formUsername" class="form-control" placeholder="name@example.com"/>
				<label class="form-label" for="formUsername">Username</label>
			</div>
		
			<!-- Password input -->
			<div class="form-floating mb-4">
				<input disabled type="password" id="formPassword" class="form-control" placeholder="gattina12"/>
				<label class="form-label" for="formPassword">Password</label>
			</div>
		
			<!-- 2 column grid layout for inline styling -->
			<div class="row mb-4">
				<div class="col d-flex justify-content-center">
					<!-- Checkbox -->
					<div class="form-check">
						<input class="form-check-input" type="checkbox" value="" id="form2Example31" checked />
						<label class="form-check-label" for="form2Example31">Remember me</label>
					</div>
				</div>
			</div>
		
			<!-- Submit button -->
			<button style="width: 100%" type="button" class="btn btn-primary btn-block mb-4" @click="login">Sign in</button>
            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
            <LoadingSpinner :loading="loading" />
		</form>
    </div>
        </div>
    </div></div>
	</div>
</template>

<style>
</style>
