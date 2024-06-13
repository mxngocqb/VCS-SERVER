import { mailService } from "@/plugins/axios/mail_sender/mailService";
import { useForm } from "vee-validate";
import * as yup from "yup";

export default function () {
    const schema = yup.object({
        email: yup.string().email().required(),
        from: yup.date().optional(),
        to: yup.date().min(yup.ref("from")).optional(),
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
    const [from, fromAttrs] = defineField("from");
    const [to, toAttrs] = defineField("to");

    const onSubmit = handleSubmit(async (value) => {
        mailService.sendReport({
            email: value["email"],
            from: value["from"] ? new Date(value["from"]).toISOString() : "",
            to: value["to"] ? new Date(value["to"]).toISOString() : "",
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
        from,
        to,
    };
}
