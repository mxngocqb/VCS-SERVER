import { userService } from "@/plugins/axios/user/userService";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        email: yup.string().email().required(),
        password: yup.string().min(6).required(),
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
    const [password, passwordAttrs] = defineField("password");
    const [fullName, fullNameAttrs] = defineField("fullName");
    const [phone, phoneAttrs] = defineField("phone");
    const [roles, rolesAttrs] = defineField("roles");

    const onSubmit = handleSubmit(async (value) => {
        userService.createUser({
            email: value["email"],
            fullName: value["fullName"],
            password: value["password"],
            roles: value["roles"] === "User" ? 2 : 1,
            phone: value["phone"],
        });
        return true;
    });
    return {
        resetForm,
        setFieldValue,
        onSubmit,
        errors,
        isValidating,
        email,
        password,
        fullName,
        phone,
        roles,
    };
}
