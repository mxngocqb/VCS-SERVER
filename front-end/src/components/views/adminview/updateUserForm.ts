import { userService } from "@/plugins/axios/user/userService";
import { useUserStore } from "@/stores/userStore";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        email: yup.string().email().required(),
        fullName: yup.string().required(),
        phone: yup.string().optional(),
        roles: yup.string().required(),
    });

    const {
        handleSubmit,
        errors,
        isValidating,
        defineField,
        setFieldValue,
        resetForm,
    } = useForm({
        validationSchema: schema,
    });

    const [email, emailAttrs] = defineField("email");
    const [fullName, fullNameAttrs] = defineField("fullName");
    const [phone, phoneAttrs] = defineField("phone");
    const [roles, rolesAttrs] = defineField("roles");
    const userStore = useUserStore();
    const onSubmit = handleSubmit(async (value) => {
        if (userStore.selectedUserComputed.value?.id) {
            return userService.updateUser(
                userStore.selectedUserComputed.value?.id,
                {
                    email: value["email"],
                    fullName: value["fullName"],
                    roles: value["roles"] === "User" ? 2 : 1,
                    phone: value["phone"],
                }
            );
        }
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        onSubmit,
        errors,
        isValidating,
        email,
        fullName,
        phone,
        roles,
    };
}
